package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"gandalf/middlewares"
	"gandalf/models"
	"gandalf/services"
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

type updateRecorder struct {
	uuid     uuid.UUID
	userData validators.UserUpdateData
}

type verificateRecorder struct {
	called bool
}

type mockUserService struct {
	createRecorder     *createRecorder
	readRecorder       *uuidRecorder
	updateRecorder     *updateRecorder
	deleteRecorder     *uuidRecorder
	softdeleteRecorder *uuidRecorder
	verificateRecorder *verificateRecorder

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

func newMockedService(createError error, readError error, updateError error, deleteError error, softdeleteError error) mockUserService {
	return mockUserService{
		createRecorder:     new(createRecorder),
		readRecorder:       new(uuidRecorder),
		updateRecorder:     new(updateRecorder),
		deleteRecorder:     new(uuidRecorder),
		softdeleteRecorder: new(uuidRecorder),
		verificateRecorder: new(verificateRecorder),
		createError:        createError,
		readError:          readError,
		updateError:        updateError,
		deleteError:        deleteError,
		softdeleteError:    softdeleteError,
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
	middleware.hasScopesCalled = true
	middleware.requestedScopes = &scopes
	return func(c *gin.Context) {}
}

func (middleware *mockAuthBearerMiddleware) GetAuthorizedUser(c *gin.Context) *models.User {
	middleware.getAuthorizedUserCalled = true
	return middleware.authorizedUser
}

func setupUserRouter(authBearerMiddleware middlewares.IAuthBearerMiddleware, userService services.IUserService) *gin.Engine {
	router := gin.Default()
	RegisterUserRoutes(router, authBearerMiddleware, userService)
	return router
}

func TestCreateUser(t *testing.T) {
	assert := require.New(t)

	t.Run("Test create user successfully", func(t *testing.T) {
		userService := newMockedService(nil, nil, nil, nil, nil)
		authBearerMiddleware := newMockAuthBearerMiddleware(nil)
		router := setupUserRouter(authBearerMiddleware, &userService)
		var response gin.H

		email := "test@test.com"
		password := "testtesttesttest"
		name := "test"
		surname := "test"
		birthdayStr := "2021-02-15T00:00:00Z"
		birthdayDat, _ := time.Parse(time.RFC3339, birthdayStr)
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

		assert.Equal(recorder.Result().StatusCode, http.StatusCreated)
		assert.Equal(userService.createRecorder.userData.Email, email)
		assert.Equal(userService.createRecorder.userData.Password, password)
		assert.Equal(userService.createRecorder.userData.Name, name)
		assert.Equal(userService.createRecorder.userData.Surname, surname)
		assert.Equal(userService.createRecorder.userData.Birthday, birthdayDat)
		assert.False(authBearerMiddleware.hasScopesCalled)
		assert.False(authBearerMiddleware.getAuthorizedUserCalled)
	})

	t.Run("Test create user binding error", func(t *testing.T) {
		userService := newMockedService(nil, nil, nil, nil, nil)
		router := setupUserRouter(newMockAuthBearerMiddleware(nil), &userService)
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
		userService := newMockedService(expectedError, nil, nil, nil, nil)
		router := setupUserRouter(newMockAuthBearerMiddleware(nil), &userService)
		var response gin.H

		email := "test@test.com"
		password := "testtesttesttest"
		name := "test"
		surname := "test"
		birthdayStr := "2021-02-15T00:00:00Z"
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

func TestVerificateUser(t *testing.T) {
	assert := require.New(t)

	t.Run("Test verificate user successfully", func(t *testing.T) {
		authorizedUser := models.NewUser(
			"test@test.com",
			"testestestestest",
			"test",
			"test",
			time.Now(),
			"+34666666666",
		)
		expectedScopes := []string{services.ScopeUserVerify}
		userService := newMockedService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupUserRouter(authMiddleware, &userService)

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("PATCH", "/users/verify", bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)

		assert.True(authMiddleware.hasScopesCalled)
		assert.True(authMiddleware.getAuthorizedUserCalled)
		assert.True(userService.verificateRecorder.called)
		assert.Equal(recorder.Result().StatusCode, http.StatusOK)
		assert.Equal(authMiddleware.requestedScopes, expectedScopes)
	})
}
