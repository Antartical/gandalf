package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"gandalf/models"
	"gandalf/security"
	"gandalf/services"
	"gandalf/tests"
	"gandalf/validators"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type authenticateRecorder struct {
	credentials validators.Credentials
}

type generateTokensRecorder struct {
	user   models.User
	scopes []string
}

type getAuthorizedUserRecorder struct {
	accessToken string
	scopes      []string
}

type refreshTokenRecorder struct {
	accessToken  string
	refreshToken string
}

type mockAuthService struct {
	authenticateRecorder      *authenticateRecorder
	generateTokensRecorder    *generateTokensRecorder
	getAuthorizedUserRecorder *getAuthorizedUserRecorder
	refreshTokenRecorder      *refreshTokenRecorder

	authenticateError      error
	getAuthorizedUserError error
	refreshTokenError      error

	returnedUser *models.User
}

func newMockedAuthService(
	returnedUser *models.User,
	authenticateError error,
	getAuthorizedUserError error,
	refreshTokenError error,
) *mockAuthService {
	return &mockAuthService{
		authenticateRecorder:      new(authenticateRecorder),
		generateTokensRecorder:    new(generateTokensRecorder),
		getAuthorizedUserRecorder: new(getAuthorizedUserRecorder),
		refreshTokenRecorder:      new(refreshTokenRecorder),
		authenticateError:         authenticateError,
		getAuthorizedUserError:    getAuthorizedUserError,
		refreshTokenError:         refreshTokenError,
		returnedUser:              returnedUser,
	}
}

func (service *mockAuthService) Authenticate(credentials validators.Credentials) (*models.User, error) {
	service.authenticateRecorder.credentials = credentials
	return service.returnedUser, service.authenticateError
}

func (service *mockAuthService) GenerateTokens(user models.User, scopes []string) services.AuthTokens {
	service.generateTokensRecorder.user = user
	service.generateTokensRecorder.scopes = scopes
	return services.AuthTokens{AccessToken: "", RefreshToken: ""}
}

func (service *mockAuthService) GetAuthorizedUser(accessToken string, scopes []string) (*models.User, error) {
	service.getAuthorizedUserRecorder.accessToken = accessToken
	service.getAuthorizedUserRecorder.scopes = scopes
	return service.returnedUser, service.getAuthorizedUserError
}

func (service *mockAuthService) RefreshToken(accessToken string, refreshToken string) (*services.AuthTokens, error) {
	service.refreshTokenRecorder.accessToken = accessToken
	service.refreshTokenRecorder.refreshToken = refreshToken
	return &services.AuthTokens{AccessToken: "", RefreshToken: refreshToken}, service.refreshTokenError
}

func setupAuthRouter(authService services.IAuthService) *gin.Engine {
	router := gin.Default()
	RegisterAuthRoutes(router, authService)
	return router
}

func TestLogin(t *testing.T) {
	assert := require.New(t)

	t.Run("Test login successfully", func(t *testing.T) {
		user := tests.UserFactory()
		expectedScopes := security.GroupUserAll
		authService := newMockedAuthService(&user, nil, nil, nil)
		router := setupAuthRouter(authService)
		var response gin.H

		payload, _ := json.Marshal(map[string]interface{}{
			"email":    user.Email,
			"password": user.Password,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusOK)
		assert.Equal(authService.authenticateRecorder.credentials.Email, user.Email)
		assert.Equal(authService.authenticateRecorder.credentials.Password, user.Password)
		assert.Equal(authService.generateTokensRecorder.user.Email, user.Email)
		assert.Equal(authService.generateTokensRecorder.scopes, expectedScopes)
	})

	t.Run("Test login wrong payload", func(t *testing.T) {
		user := tests.UserFactory()
		authService := newMockedAuthService(nil, nil, nil, nil)
		router := setupAuthRouter(authService)

		payload, _ := json.Marshal(map[string]string{
			"email": user.Email,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("Test login forbidden user", func(t *testing.T) {
		raisedError := errors.New("wrong")
		user := tests.UserFactory()
		authService := newMockedAuthService(nil, raisedError, nil, nil)
		router := setupAuthRouter(authService)

		payload, _ := json.Marshal(map[string]interface{}{
			"email":    user.Email,
			"password": user.Password,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(recorder.Result().StatusCode, http.StatusForbidden)
	})
}

func TestRefresh(t *testing.T) {
	assert := require.New(t)

	t.Run("Test refresh token succesfully", func(t *testing.T) {
		authService := newMockedAuthService(nil, nil, nil, nil)
		router := setupAuthRouter(authService)
		accessToken := "testaccess"
		refreshToken := "testrefresh"
		var response gin.H

		payload, _ := json.Marshal(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusOK)
		assert.Equal(authService.refreshTokenRecorder.accessToken, accessToken)
		assert.Equal(authService.refreshTokenRecorder.refreshToken, refreshToken)
	})

	t.Run("Test refresh token wrong payload", func(t *testing.T) {
		authService := newMockedAuthService(nil, nil, nil, nil)
		router := setupAuthRouter(authService)
		accessToken := "testaccess"

		payload, _ := json.Marshal(map[string]string{
			"access_token": accessToken,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("Test refresh token unrelated tokens", func(t *testing.T) {
		raisedError := errors.New("wrong")
		authService := newMockedAuthService(nil, nil, nil, raisedError)
		router := setupAuthRouter(authService)
		accessToken := "testaccess"
		refreshToken := "testrefresh"
		var response gin.H

		payload, _ := json.Marshal(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})
}
