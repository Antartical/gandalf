package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"gandalf/middlewares"
	"gandalf/models"
	"gandalf/security"
	"gandalf/services"
	"gandalf/tests"
	"gandalf/validators"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

type createAppRecorder struct {
	appData validators.AppCreateData
	user    models.User
}

type readAppRecorder struct {
	uuid uuid.UUID
}

type readByClientAppRecorder struct {
	clientId uuid.UUID
}

type updateAppRecorder struct {
	uuid          uuid.UUID
	appUpdateData validators.AppUpdateData
}

type deleteAppRecorder struct {
	uuid uuid.UUID
}

type mockAppService struct {
	createAppRecorder       *createAppRecorder
	readAppRecorder         *readAppRecorder
	readByClientAppRecorder *readByClientAppRecorder
	updateAppRecorder       *updateAppRecorder
	deleteAppRecorder       *deleteAppRecorder

	createError       error
	readError         error
	readByClientError error
	updateError       error
	deleteError       error
}

func (service *mockAppService) Create(appData validators.AppCreateData, user models.User) (*models.App, error) {
	*service.createAppRecorder = createAppRecorder{appData, user}
	return &models.App{}, service.createError
}

func (service *mockAppService) Read(uuid uuid.UUID) (*models.App, error) {
	*service.readAppRecorder = readAppRecorder{uuid}
	return &models.App{}, service.readError
}

func (service *mockAppService) ReadByClientID(clientID uuid.UUID) (*models.App, error) {
	*service.readByClientAppRecorder = readByClientAppRecorder{clientID}
	return &models.App{}, service.readByClientError
}

func (service *mockAppService) Update(uuid uuid.UUID, appData validators.AppUpdateData) (*models.App, error) {
	*service.updateAppRecorder = updateAppRecorder{uuid, appData}
	return &models.App{}, service.updateError
}

func (service *mockAppService) Delete(uuid uuid.UUID) error {
	*service.deleteAppRecorder = deleteAppRecorder{uuid}
	return service.deleteError
}

func newMockedAppService(createError error, readError error, readByClientError error, updateError error, deleteError error) mockAppService {
	return mockAppService{
		createAppRecorder:       new(createAppRecorder),
		readAppRecorder:         new(readAppRecorder),
		readByClientAppRecorder: new(readByClientAppRecorder),
		updateAppRecorder:       new(updateAppRecorder),
		deleteAppRecorder:       new(deleteAppRecorder),
		createError:             createError,
		readError:               readError,
		readByClientError:       readByClientError,
		updateError:             updateError,
		deleteError:             deleteError,
	}
}

func setupOauth2Router(
	authBearerMiddleware middlewares.IAuthBearerMiddleware,
	authService services.IAuthService,
	userService services.IUserService,
	appService services.IAppService,
) *gin.Engine {
	router := gin.Default()
	RegisterOauth2Routes(
		router, authBearerMiddleware,
		authService, userService, appService,
	)
	return router
}

func TestOauth2Login(t *testing.T) {
	assert := require.New(t)

	t.Run("Test oauth2 login success", func(t *testing.T) {
		user := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authBearerMiddleware := newMockAuthBearerMiddleware(nil)
		authService := newMockedAuthService(&user, nil, nil, nil, nil, nil)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupOauth2Router(
			authBearerMiddleware,
			authService,
			&userService,
			&appService,
		)

		expectedScopes := security.GroupUserOauth2Request
		var response gin.H
		payload, _ := json.Marshal(map[string]interface{}{
			"email":    user.Email,
			"password": user.Password,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/oauth/login", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(http.StatusOK, recorder.Result().StatusCode)
		assert.Equal(authService.authenticateRecorder.credentials.Email, user.Email)
		assert.Equal(authService.authenticateRecorder.credentials.Password, user.Password)
		assert.Equal(authService.generateTokensRecorder.user.Email, user.Email)
		assert.Equal(authService.generateTokensRecorder.scopes, expectedScopes)
	})

	t.Run("Test oauth2 login success bad request", func(t *testing.T) {
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authBearerMiddleware := newMockAuthBearerMiddleware(nil)
		authService := newMockedAuthService(nil, nil, nil, nil, nil, nil)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupOauth2Router(
			authBearerMiddleware,
			authService,
			&userService,
			&appService,
		)

		user := tests.UserFactory()
		payload, _ := json.Marshal(map[string]interface{}{
			"email": user.Email,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/oauth/login", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("Test oauth2 login authorization error", func(t *testing.T) {
		raisedError := errors.New("wrong")

		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authBearerMiddleware := newMockAuthBearerMiddleware(nil)
		authService := newMockedAuthService(nil, raisedError, nil, nil, nil, nil)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupOauth2Router(
			authBearerMiddleware,
			authService,
			&userService,
			&appService,
		)

		user := tests.UserFactory()
		payload, _ := json.Marshal(map[string]interface{}{
			"email":    user.Email,
			"password": user.Password,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/oauth/login", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(recorder.Result().StatusCode, http.StatusForbidden)
	})

}

func TestOauth2Authorize(t *testing.T) {
	assert := require.New(t)

	t.Run("Test oauth2 authorize success", func(t *testing.T) {
		user := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		authService := newMockedAuthService(&user, nil, nil, nil, nil, nil)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupOauth2Router(
			authBearerMiddleware,
			authService,
			&userService,
			&appService,
		)

		var response gin.H
		uuid, _ := uuid.NewV4()
		redirectUrl := faker.Internet().Url()
		state := faker.RandomString(10)
		payload, _ := json.Marshal(map[string]interface{}{
			"client_id":    uuid,
			"redirect_uri": redirectUrl,
			"scopes":       security.GroupAdmin,
			"state":        state,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/oauth/authorize", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Contains(recorder.Result().Header["Location"][0], redirectUrl)
		assert.Equal(http.StatusFound, recorder.Result().StatusCode)
	})

	t.Run("Test oauth2 authorize bad request payload", func(t *testing.T) {
		user := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		authService := newMockedAuthService(&user, nil, nil, nil, nil, nil)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupOauth2Router(
			authBearerMiddleware,
			authService,
			&userService,
			&appService,
		)

		uuid, _ := uuid.NewV4()
		redirectUrl := faker.Internet().Url()
		state := faker.RandomString(10)
		payload, _ := json.Marshal(map[string]interface{}{
			"client_id":    uuid,
			"redirect_uri": redirectUrl,
			"scopis":       security.GroupAdmin,
			"state":        state,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/oauth/authorize", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(http.StatusBadRequest, recorder.Result().StatusCode)
	})

	t.Run("Test oauth2 authorize bad request app does not exist", func(t *testing.T) {
		expectedError := errors.New("Whoops!")
		user := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		authService := newMockedAuthService(&user, nil, nil, nil, nil, nil)
		appService := newMockedAppService(nil, nil, expectedError, nil, nil)
		router := setupOauth2Router(
			authBearerMiddleware,
			authService,
			&userService,
			&appService,
		)

		uuid, _ := uuid.NewV4()
		redirectUrl := faker.Internet().Url()
		state := faker.RandomString(10)
		payload, _ := json.Marshal(map[string]interface{}{
			"client_id":    uuid,
			"redirect_uri": redirectUrl,
			"scopes":       security.GroupAdmin,
			"state":        state,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/oauth/authorize", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(http.StatusBadRequest, recorder.Result().StatusCode)
	})

	t.Run("Test oauth2 authorize error", func(t *testing.T) {
		expectedError := errors.New("Whoops!")
		user := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		authService := newMockedAuthService(&user, nil, nil, nil, expectedError, nil)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupOauth2Router(
			authBearerMiddleware,
			authService,
			&userService,
			&appService,
		)

		uuid, _ := uuid.NewV4()
		redirectUrl := faker.Internet().Url()
		state := faker.RandomString(10)
		payload, _ := json.Marshal(map[string]interface{}{
			"client_id":    uuid,
			"redirect_uri": redirectUrl,
			"scopes":       security.GroupAdmin,
			"state":        state,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/oauth/authorize", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(http.StatusBadRequest, recorder.Result().StatusCode)
	})

}

func TestOauth2Token(t *testing.T) {
	assert := require.New(t)

	t.Run("Test oauth2 token success", func(t *testing.T) {
		user := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		authService := newMockedAuthService(&user, nil, nil, nil, nil, nil)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupOauth2Router(
			authBearerMiddleware,
			authService,
			&userService,
			&appService,
		)

		var response gin.H
		uuid, _ := uuid.NewV4()
		payload, _ := json.Marshal(map[string]interface{}{
			"grant_type":    "authorization_code",
			"client_id":     uuid,
			"client_secret": faker.RandomString(10),
			"code":          faker.RandomString(10),
			"redirect_uri":  faker.Internet().Url(),
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(http.StatusOK, recorder.Result().StatusCode)
	})

	t.Run("Test oauth2 wrong payload", func(t *testing.T) {
		user := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		authService := newMockedAuthService(&user, nil, nil, nil, nil, nil)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupOauth2Router(
			authBearerMiddleware,
			authService,
			&userService,
			&appService,
		)

		uuid, _ := uuid.NewV4()
		payload, _ := json.Marshal(map[string]interface{}{
			"grant_type":    "token",
			"client_id":     uuid,
			"client_secret": faker.RandomString(10),
			"code":          faker.RandomString(10),
			"redirect_uri":  faker.Internet().Url(),
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(http.StatusBadRequest, recorder.Result().StatusCode)
	})

	t.Run("Test oauth2 app error", func(t *testing.T) {
		expectedError := errors.New("Whoops!")
		user := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		authService := newMockedAuthService(&user, nil, nil, nil, nil, nil)
		appService := newMockedAppService(nil, nil, expectedError, nil, nil)
		router := setupOauth2Router(
			authBearerMiddleware,
			authService,
			&userService,
			&appService,
		)

		uuid, _ := uuid.NewV4()
		payload, _ := json.Marshal(map[string]interface{}{
			"grant_type":    "authorization_code",
			"client_id":     uuid,
			"client_secret": faker.RandomString(10),
			"code":          faker.RandomString(10),
			"redirect_uri":  faker.Internet().Url(),
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(http.StatusBadRequest, recorder.Result().StatusCode)
	})

	t.Run("Test oauth2 exchange error", func(t *testing.T) {
		expectedError := errors.New("Whoops!")
		user := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		authService := newMockedAuthService(&user, nil, nil, nil, nil, expectedError)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupOauth2Router(
			authBearerMiddleware,
			authService,
			&userService,
			&appService,
		)

		uuid, _ := uuid.NewV4()
		payload, _ := json.Marshal(map[string]interface{}{
			"grant_type":    "authorization_code",
			"client_id":     uuid,
			"client_secret": faker.RandomString(10),
			"code":          faker.RandomString(10),
			"redirect_uri":  faker.Internet().Url(),
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(http.StatusUnauthorized, recorder.Result().StatusCode)
	})
}
