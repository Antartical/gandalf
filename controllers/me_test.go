package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gandalf/middlewares"
	"gandalf/security"
	"gandalf/serializers"
	"gandalf/services"
	"gandalf/tests"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func setupMeRouter(
	authBearerMiddleware middlewares.IAuthBearerMiddleware,
	authService services.IAuthService,
	userService services.IUserService,
	appService services.IAppService,
	pelipperService services.IPelipperService,
) *gin.Engine {
	router := gin.Default()
	RegisterMeRoutes(
		router, authBearerMiddleware,
		authService, userService,
		appService, pelipperService,
	)
	return router
}

func TestReadMe(t *testing.T) {
	assert := require.New(t)

	t.Run("Test Me", func(t *testing.T) {
		authorizedUser := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService,
			&appService,
			newPelipperServiceMock(),
		)
		var response gin.H

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("GET", "/me", bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusOK)
		assert.True(authMiddleware.getAuthorizedUserCalled)
	})
}

func TestUpdateMe(t *testing.T) {
	assert := require.New(t)

	t.Run("Test update me successfully", func(t *testing.T) {
		authorizedUser := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService,
			&appService,
			newPelipperServiceMock(),
		)
		var response gin.H

		password := faker.Internet().Password(10, 14)
		phone := fmt.Sprintf(
			"+%s%s",
			faker.PhoneNumber().AreaCode(),
			faker.PhoneNumber().SubscriberNumber(9),
		)
		payload, _ := json.Marshal(map[string]string{
			"password": password,
			"phone":    phone,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("PATCH", "/me", bytes.NewBuffer(payload))

		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(http.StatusOK, recorder.Result().StatusCode)
		assert.True(authMiddleware.getAuthorizedUserCalled)
		assert.Equal(userService.updateRecorder.userData.Password, password)
		assert.Equal(userService.updateRecorder.userData.Phone, phone)
	})

	t.Run("Test update me wrong payload", func(t *testing.T) {
		authorizedUser := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService,
			&appService,
			newPelipperServiceMock(),
		)
		var response gin.H

		password := faker.Internet().Password(10, 14)
		phone := faker.PhoneNumber().CellPhone()
		payload, _ := json.Marshal(map[string]string{
			"password": password,
			"phone":    phone,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("PATCH", "/me", bytes.NewBuffer(payload))

		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(http.StatusBadRequest, recorder.Result().StatusCode)
	})

	t.Run("Test update me db error", func(t *testing.T) {
		expectedError := errors.New("Whoops!")
		authorizedUser := tests.UserFactory()
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		userService := newMockedUserService(nil, nil, expectedError, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService,
			&appService,
			newPelipperServiceMock(),
		)
		var response gin.H

		password := faker.Internet().Password(10, 14)
		phone := fmt.Sprintf(
			"+%s%s",
			faker.PhoneNumber().AreaCode(),
			faker.PhoneNumber().SubscriberNumber(9),
		)
		payload, _ := json.Marshal(map[string]string{
			"password": password,
			"phone":    phone,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("PATCH", "/me", bytes.NewBuffer(payload))

		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)
		assert.Equal(http.StatusBadRequest, recorder.Result().StatusCode)
	})
}

func TestDeleteMe(t *testing.T) {
	assert := require.New(t)

	t.Run("Test Delete me successfully", func(t *testing.T) {
		authorizedUser := tests.UserFactory()
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, &appService, newPelipperServiceMock(),
		)
		var response gin.H

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("DELETE", "/me", bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusNoContent)
		assert.True(authMiddleware.getAuthorizedUserCalled)
	})

	t.Run("Test Delete db error", func(t *testing.T) {
		expectedError := errors.New("Whoops!")
		authorizedUser := tests.UserFactory()
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		userService := newMockedUserService(nil, nil, nil, expectedError, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, &appService, newPelipperServiceMock(),
		)
		var response gin.H

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("DELETE", "/me", bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
		assert.True(authMiddleware.getAuthorizedUserCalled)
	})
}

func TestVerificateMe(t *testing.T) {
	assert := require.New(t)

	t.Run("Test verificate me successfully", func(t *testing.T) {
		authorizedUser := tests.UserFactory()
		expectedScopes := []string{security.ScopeUserVerify}
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, &appService, newPelipperServiceMock(),
		)

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/me/verify", bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)

		assert.True(authMiddleware.hasScopesCalled)
		assert.True(authMiddleware.getAuthorizedUserCalled)
		assert.True(userService.verificateRecorder.called)
		assert.Equal(recorder.Result().StatusCode, http.StatusNoContent)
		assert.Equal(authMiddleware.requestedScopes, &expectedScopes)
	})
}

func TestResetMyPassword(t *testing.T) {
	assert := require.New(t)

	t.Run("Test reset my password successfully", func(t *testing.T) {
		password := "testestestestest"
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		authorizedUser := tests.UserFactory()
		payload, _ := json.Marshal(map[string]string{
			"password": password,
		})
		expectedScopes := []string{security.ScopeUserChangePassword}
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, &appService, newPelipperServiceMock(),
		)

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/me/reset-password", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.True(authMiddleware.hasScopesCalled)
		assert.True(authMiddleware.getAuthorizedUserCalled)
		assert.Equal(userService.resetPasswordRecorder.password, password)
		assert.Equal(recorder.Result().StatusCode, http.StatusNoContent)
		assert.Equal(authMiddleware.requestedScopes, &expectedScopes)
	})

	t.Run("Test reset my password wrong payload", func(t *testing.T) {
		password := "testestestestest"
		authorizedUser := tests.UserFactory()
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		payload, _ := json.Marshal(map[string]string{
			"wrong": password,
		})
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService,
			&appService,
			newPelipperServiceMock(),
		)

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/me/reset-password", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})
}

func TestGetMyApps(t *testing.T) {
	assert := require.New(t)

	t.Run("Test get my apps success", func(t *testing.T) {
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		authorizedUser := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, &appService, newPelipperServiceMock(),
		)

		var response serializers.PaginatedAppsSerializer
		recorder := httptest.NewRecorder()
		page := 3
		limit := 5
		url := fmt.Sprintf("/me/apps?page=%d&limit=%d", page, limit)
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusOK)
		assert.Equal(response.Meta.Cursor.Data.ActualPage, page)
	})

	t.Run("Test get my apps bad request", func(t *testing.T) {
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		authorizedUser := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, &appService, newPelipperServiceMock(),
		)

		var response serializers.PaginatedAppsSerializer
		recorder := httptest.NewRecorder()
		page := 3
		limit := 10000
		url := fmt.Sprintf("/me/apps?page=%d&limit=%d", page, limit)
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})
}

func TestGetMyConnectedApps(t *testing.T) {
	assert := require.New(t)

	t.Run("Test get my apps success", func(t *testing.T) {
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		authorizedUser := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, &appService, newPelipperServiceMock(),
		)

		var response serializers.PaginatedAppsPublicSerializer
		recorder := httptest.NewRecorder()
		page := 3
		limit := 5
		url := fmt.Sprintf("/me/connected-apps?page=%d&limit=%d", page, limit)
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusOK)
		assert.Equal(response.Meta.Cursor.Data.ActualPage, page)
	})

	t.Run("Test get my apps bad request", func(t *testing.T) {
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		authorizedUser := tests.UserFactory()
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		authMiddleware := newMockAuthBearerMiddleware(&authorizedUser)
		router := setupMeRouter(
			authMiddleware,
			newMockedAuthService(nil, nil, nil, nil, nil, nil),
			&userService, &appService, newPelipperServiceMock(),
		)

		var response serializers.PaginatedAppsSerializer
		recorder := httptest.NewRecorder()
		page := 3
		limit := 10000
		url := fmt.Sprintf("/me/connected-apps?page=%d&limit=%d", page, limit)
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})
}
