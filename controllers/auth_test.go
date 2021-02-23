package controllers

import (
	"gandalf/models"
	"gandalf/services"
	"gandalf/validators"
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
	return &services.AuthTokens{AccessToken: "", RefreshToken: ""}, service.refreshTokenError
}
