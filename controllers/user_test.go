package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gandalf/bindings"
	"gandalf/middlewares"
	"gandalf/models"
	"gandalf/services"
	"gandalf/tests"
	"gandalf/validators"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

type createRecorder struct {
	userData validators.UserCreateData
}

type uuidRecorder struct {
	uuid uuid.UUID
}

type emailRecorder struct {
	email string
}

type updateRecorder struct {
	uuid     uuid.UUID
	userData validators.UserUpdateData
}

type verificateRecorder struct {
	called bool
}

type resetPasswordRecorder struct {
	password string
}

type mockUserService struct {
	createRecorder        *createRecorder
	readRecorder          *uuidRecorder
	readByEmailRecorder   *emailRecorder
	updateRecorder        *updateRecorder
	deleteRecorder        *uuidRecorder
	softdeleteRecorder    *uuidRecorder
	verificateRecorder    *verificateRecorder
	resetPasswordRecorder *resetPasswordRecorder

	createError     error
	readError       error
	updateError     error
	deleteError     error
	softdeleteError error
}

func (service *mockUserService) Create(userData validators.UserCreateData) (*models.User, error) {
	*service.createRecorder = createRecorder{userData: userData}
	return &models.User{}, service.createError
}

func (service *mockUserService) Read(uuid uuid.UUID) (*models.User, error) {
	*service.readRecorder = uuidRecorder{uuid: uuid}
	return &models.User{}, service.readError
}

func (service *mockUserService) ReadByEmail(email string) (*models.User, error) {
	*service.readByEmailRecorder = emailRecorder{email: email}
	return &models.User{Email: email}, service.readError
}

func (service *mockUserService) Update(uuid uuid.UUID, userData validators.UserUpdateData) (*models.User, error) {
	*service.updateRecorder = updateRecorder{uuid: uuid, userData: userData}
	return &models.User{}, service.updateError
}

func (service *mockUserService) Delete(uuid uuid.UUID) error {
	*service.deleteRecorder = uuidRecorder{uuid: uuid}
	return service.deleteError
}

func (service *mockUserService) SoftDelete(uuid uuid.UUID) error {
	*service.softdeleteRecorder = uuidRecorder{uuid: uuid}
	return service.softdeleteError
}

func (service *mockUserService) Verificate(*models.User) {
	*service.verificateRecorder = verificateRecorder{called: true}
}

func (service *mockUserService) ResetPassword(user *models.User, password string) {
	*service.resetPasswordRecorder = resetPasswordRecorder{password}
}

func newMockedUserService(createError error, readError error, updateError error, deleteError error, softdeleteError error) mockUserService {
	return mockUserService{
		createRecorder:        new(createRecorder),
		readByEmailRecorder:   new(emailRecorder),
		readRecorder:          new(uuidRecorder),
		updateRecorder:        new(updateRecorder),
		deleteRecorder:        new(uuidRecorder),
		softdeleteRecorder:    new(uuidRecorder),
		verificateRecorder:    new(verificateRecorder),
		resetPasswordRecorder: new(resetPasswordRecorder),
		createError:           createError,
		readError:             readError,
		updateError:           updateError,
		deleteError:           deleteError,
		softdeleteError:       softdeleteError,
	}
}

type mockAuthBearerMiddleware struct {
	hasScopesCalled         bool
	requestedScopes         *[]string
	getAuthorizedUserCalled bool

	authorizedUser *models.User
}

func newMockAuthBearerMiddleware(authorizedUser *models.User) *mockAuthBearerMiddleware {
	return &mockAuthBearerMiddleware{false, new([]string), false, authorizedUser}
}

func (middleware *mockAuthBearerMiddleware) HasScopes(scopes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.hasScopesCalled = true
		middleware.requestedScopes = &scopes
	}
}

func (middleware *mockAuthBearerMiddleware) GetAuthorizedUser(c *gin.Context) *models.User {
	middleware.getAuthorizedUserCalled = true
	return middleware.authorizedUser
}

func setupUserRouter(
	authBearerMiddleware middlewares.IAuthBearerMiddleware,
	authService services.IAuthService,
	userService services.IUserService,
	pelipperService services.IPelipperService,
) *gin.Engine {
	router := gin.Default()
	RegisterUserRoutes(
		router, authBearerMiddleware,
		authService, userService, pelipperService,
	)
	return router
}

type sendUserVerifyEmailRecorder struct {
	data validators.PelipperUserVerifyEmail
}

type sendUserChangePasswordEmailRecorder struct {
	data validators.PelipperUserChangePassword
}
type pelipperServiceMock struct {
	sendUserVerifyEmailRecorder         *sendUserVerifyEmailRecorder
	sendUserChangePasswordEmailRecorder *sendUserChangePasswordEmailRecorder
}

func newPelipperServiceMock() *pelipperServiceMock {
	return &pelipperServiceMock{
		sendUserVerifyEmailRecorder:         new(sendUserVerifyEmailRecorder),
		sendUserChangePasswordEmailRecorder: new(sendUserChangePasswordEmailRecorder),
	}
}

func (service *pelipperServiceMock) SendUserVerifyEmail(data validators.PelipperUserVerifyEmail) {
	service.sendUserVerifyEmailRecorder.data = data
}

func (service *pelipperServiceMock) SendUserChangePasswordEmail(data validators.PelipperUserChangePassword) {
	service.sendUserChangePasswordEmailRecorder.data = data
}

func TestCreateUser(t *testing.T) {
	assert := require.New(t)

	t.Run("Test create user successfully", func(t *testing.T) {
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		pelipperService := newPelipperServiceMock()
		authBearerMiddleware := newMockAuthBearerMiddleware(nil)
		router := setupUserRouter(
			authBearerMiddleware, newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, pelipperService,
		)
		var response gin.H

		email := "test@test.com"
		password := "testtesttesttest"
		name := "test"
		surname := "test"
		birthdayStr := "2021-02-15T00:00:00Z"
		birthdayDat, _ := time.Parse(time.RFC3339, birthdayStr)
		birthday := bindings.BirthDate(birthdayDat)
		payload, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
			"name":     name,
			"Surname":  surname,
			"Birthday": "2021-02-15",
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusCreated)
		assert.Equal(userService.createRecorder.userData.Email, email)
		assert.Equal(userService.createRecorder.userData.Password, password)
		assert.Equal(userService.createRecorder.userData.Name, name)
		assert.Equal(userService.createRecorder.userData.Surname, surname)
		assert.Equal(userService.createRecorder.userData.Birthday, birthday)
		assert.False(authBearerMiddleware.hasScopesCalled)
		assert.False(authBearerMiddleware.getAuthorizedUserCalled)
	})

	t.Run("Test create user binding error", func(t *testing.T) {
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		router := setupUserRouter(
			newMockAuthBearerMiddleware(nil),
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, newPelipperServiceMock(),
		)
		var response gin.H

		payload, _ := json.Marshal(map[string]string{
			"email": "wrong email",
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("Test create user service error", func(t *testing.T) {
		expectedError := errors.New("create error")
		userService := newMockedUserService(expectedError, nil, nil, nil, nil)
		router := setupUserRouter(
			newMockAuthBearerMiddleware(nil),
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, newPelipperServiceMock(),
		)
		var response gin.H

		email := "test@test.com"
		password := "testtesttesttest"
		name := "test"
		surname := "test"
		birthdayStr := "2021-02-15"
		payload, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
			"name":     name,
			"Surname":  surname,
			"Birthday": birthdayStr,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
		assert.Equal(response["error"], expectedError.Error())
	})
}

func TestReadUser(t *testing.T) {
	assert := require.New(t)

	t.Run("Test Read successfully", func(t *testing.T) {
		authorizedUser := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupUserRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, newPelipperServiceMock(),
		)
		var response gin.H

		recorder := httptest.NewRecorder()
		uuid, _ := uuid.NewV4()
		url := fmt.Sprintf("/users/%s", uuid.String())
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(http.StatusOK, recorder.Result().StatusCode)
	})

	t.Run("Test Read wrong uuid param", func(t *testing.T) {
		authorizedUser := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupUserRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, newPelipperServiceMock(),
		)
		var response gin.H

		recorder := httptest.NewRecorder()
		url := fmt.Sprintf("/users/invent")
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("Test Read not found", func(t *testing.T) {
		expectedError := errors.New("Whoops!")
		authorizedUser := tests.UserFactory()
		userService := newMockedUserService(nil, expectedError, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupUserRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, newPelipperServiceMock(),
		)
		var response gin.H

		recorder := httptest.NewRecorder()
		uuid, _ := uuid.NewV4()
		url := fmt.Sprintf("/users/%s", uuid)
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusNotFound)
	})
}
