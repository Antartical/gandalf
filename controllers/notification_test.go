package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"gandalf/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupNotificationRouter(
	authService services.IAuthService,
	userService services.IUserService,
	pelipperService services.IPelipperService,
) *gin.Engine {
	router := gin.Default()
	RegisterNotificationRoutes(
		router, authService,
		userService, pelipperService,
	)
	return router
}

func TestUserResendVerificationEmail(t *testing.T) {
	assert := require.New(t)

	t.Run("Test resend verification email successfully", func(t *testing.T) {
		authService := newMockedAuthService(nil, nil, nil, nil, nil, nil)
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		pelipperService := newPelipperServiceMock()
		router := setupNotificationRouter(
			authService, &userService, pelipperService,
		)
		var response gin.H
		email := "test@test.com"

		payload, _ := json.Marshal(map[string]string{
			"email": email,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/notifications/emails/verify", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusNoContent)
		assert.Equal(userService.readByEmailRecorder.email, email)
		assert.Equal(authService.generateTokensRecorder.user.Email, email)
	})

	t.Run("Test resend verification email wrong payload", func(t *testing.T) {
		authService := newMockedAuthService(nil, nil, nil, nil, nil, nil)
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		pelipperService := newPelipperServiceMock()
		router := setupNotificationRouter(
			authService, &userService, pelipperService,
		)
		var response gin.H
		email := "test@test.com"

		payload, _ := json.Marshal(map[string]string{
			"wrong": email,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/notifications/emails/verify", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("Test resend verification email not registered", func(t *testing.T) {
		authService := newMockedAuthService(nil, nil, nil, nil, nil, nil)
		userService := newMockedUserService(nil, errors.New("not found"), nil, nil, nil)
		pelipperService := newPelipperServiceMock()
		router := setupNotificationRouter(
			authService, &userService, pelipperService,
		)
		var response gin.H
		email := "test@test.com"

		payload, _ := json.Marshal(map[string]string{
			"email": email,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/notifications/emails/verify", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusCreated)
		assert.Equal(userService.readByEmailRecorder.email, email)
	})
}

func TestUserResendResetPasswordEmail(t *testing.T) {
	assert := require.New(t)

	t.Run("Test resend change password email successfully", func(t *testing.T) {
		authService := newMockedAuthService(nil, nil, nil, nil, nil, nil)
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		pelipperService := newPelipperServiceMock()
		router := setupNotificationRouter(
			authService, &userService, pelipperService,
		)
		var response gin.H
		email := "test@test.com"

		payload, _ := json.Marshal(map[string]string{
			"email": email,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/notifications/emails/reset-password", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusNoContent)
		assert.Equal(userService.readByEmailRecorder.email, email)
		assert.Equal(authService.generateTokensRecorder.user.Email, email)
	})

	t.Run("Test resend change password email wrong payload", func(t *testing.T) {
		authService := newMockedAuthService(nil, nil, nil, nil, nil, nil)
		userService := newMockedUserService(nil, nil, nil, nil, nil)
		pelipperService := newPelipperServiceMock()
		router := setupNotificationRouter(
			authService, &userService, pelipperService,
		)
		var response gin.H
		email := "test@test.com"

		payload, _ := json.Marshal(map[string]string{
			"wrong": email,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/notifications/emails/reset-password", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("Test resend change password email not registered", func(t *testing.T) {
		authService := newMockedAuthService(nil, nil, nil, nil, nil, nil)
		userService := newMockedUserService(nil, errors.New("not found"), nil, nil, nil)
		pelipperService := newPelipperServiceMock()
		router := setupNotificationRouter(
			authService, &userService, pelipperService,
		)
		var response gin.H
		email := "test@test.com"

		payload, _ := json.Marshal(map[string]string{
			"email": email,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/notifications/emails/reset-password", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusCreated)
		assert.Equal(userService.readByEmailRecorder.email, email)
	})
}
