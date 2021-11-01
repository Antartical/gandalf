package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gandalf/middlewares"
	"gandalf/services"
	"gandalf/tests"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func setupAppRouter(
	authBearerMiddleware middlewares.IAuthBearerMiddleware,
	appService services.IAppService,
) *gin.Engine {
	router := gin.Default()
	RegisterAppRoutes(
		router, authBearerMiddleware, appService,
	)
	return router
}

func TestCreateApp(t *testing.T) {
	assert := require.New(t)

	t.Run("Test create app successfully", func(t *testing.T) {
		user := tests.UserFactory()
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupAppRouter(authBearerMiddleware, &appService)

		name := faker.Company().Name()
		iconUrl := faker.Internet().Url()
		redirectUrls := []string{faker.Internet().Url()}

		payload, _ := json.Marshal(map[string]interface{}{
			"name":          name,
			"icon_url":      iconUrl,
			"redirect_urls": redirectUrls,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/apps", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(recorder.Result().StatusCode, http.StatusCreated)
		assert.Equal(appService.createAppRecorder.appData.Name, name)
		assert.Equal(appService.createAppRecorder.appData.IconUrl, iconUrl)
		assert.Equal(appService.createAppRecorder.appData.RedirectUrls, redirectUrls)
		assert.True(authBearerMiddleware.getAuthorizedUserCalled)
	})

	t.Run("Test create app wrong payload", func(t *testing.T) {
		user := tests.UserFactory()
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupAppRouter(authBearerMiddleware, &appService)

		payload, _ := json.Marshal(map[string]interface{}{
			"fakeparam": "FAkE",
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/apps", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("Test create app db error", func(t *testing.T) {
		expectedError := errors.New("Whoops!")
		user := tests.UserFactory()
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		appService := newMockedAppService(expectedError, nil, nil, nil, nil)
		router := setupAppRouter(authBearerMiddleware, &appService)

		name := faker.Company().Name()
		iconUrl := faker.Internet().Url()
		redirectUrls := []string{faker.Internet().Url()}

		payload, _ := json.Marshal(map[string]interface{}{
			"name":          name,
			"icon_url":      iconUrl,
			"redirect_urls": redirectUrls,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/apps", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})
}

func TestReadApp(t *testing.T) {
	assert := require.New(t)

	t.Run("Test read app successfully", func(t *testing.T) {
		user := tests.UserFactory()
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupAppRouter(authBearerMiddleware, &appService)

		recorder := httptest.NewRecorder()
		uuid, _ := uuid.NewV4()
		url := fmt.Sprintf("/apps/%s", uuid.String())
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)

		assert.Equal(http.StatusOK, recorder.Result().StatusCode)
	})

	t.Run("Test read app wrong payload", func(t *testing.T) {
		user := tests.UserFactory()
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupAppRouter(authBearerMiddleware, &appService)

		recorder := httptest.NewRecorder()
		url := "/apps/invent"
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)

		assert.Equal(http.StatusBadRequest, recorder.Result().StatusCode)
	})

	t.Run("Test read app db error", func(t *testing.T) {
		expectedError := errors.New("Whoops!")
		user := tests.UserFactory()
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		appService := newMockedAppService(nil, expectedError, nil, nil, nil)
		router := setupAppRouter(authBearerMiddleware, &appService)

		recorder := httptest.NewRecorder()
		uuid, _ := uuid.NewV4()
		url := fmt.Sprintf("/apps/%s", uuid.String())
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)

		assert.Equal(http.StatusNotFound, recorder.Result().StatusCode)
	})
}

func TestReadAppByClientID(t *testing.T) {
	assert := require.New(t)

	t.Run("Test read app successfully", func(t *testing.T) {
		user := tests.UserFactory()
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupAppRouter(authBearerMiddleware, &appService)

		recorder := httptest.NewRecorder()
		uuid, _ := uuid.NewV4()
		url := fmt.Sprintf("/apps/public/%s", uuid.String())
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)

		assert.Equal(http.StatusOK, recorder.Result().StatusCode)
	})

	t.Run("Test read app wrong payload", func(t *testing.T) {
		user := tests.UserFactory()
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		appService := newMockedAppService(nil, nil, nil, nil, nil)
		router := setupAppRouter(authBearerMiddleware, &appService)

		recorder := httptest.NewRecorder()
		url := "/apps/public/invent"
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)

		assert.Equal(http.StatusBadRequest, recorder.Result().StatusCode)
	})

	t.Run("Test read app db error", func(t *testing.T) {
		expectedError := errors.New("Whoops!")
		user := tests.UserFactory()
		authBearerMiddleware := newMockAuthBearerMiddleware(&user)
		appService := newMockedAppService(nil, nil, expectedError, nil, nil)
		router := setupAppRouter(authBearerMiddleware, &appService)

		recorder := httptest.NewRecorder()
		uuid, _ := uuid.NewV4()
		url := fmt.Sprintf("/apps/public/%s", uuid.String())
		request, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		router.ServeHTTP(recorder, request)

		assert.Equal(http.StatusNotFound, recorder.Result().StatusCode)
	})
}
