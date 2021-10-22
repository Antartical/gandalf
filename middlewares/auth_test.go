package middlewares

import (
	"bytes"
	"fmt"
	"gandalf/models"
	"gandalf/services"
	"gandalf/tests"
	"gandalf/validators"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type getAuthorizedUserRecorder struct {
	accessToken string
	scopes      []string
}
type authServiceMock struct {
	recorder               *getAuthorizedUserRecorder
	userGetAuthorizedUser  *models.User
	errorGetAuthorizedUser error
}

func newAuthServiceMock(userGetAuthorizedUser *models.User, errorGetAuthorizedUser error) *authServiceMock {
	return &authServiceMock{
		recorder:               new(getAuthorizedUserRecorder),
		userGetAuthorizedUser:  userGetAuthorizedUser,
		errorGetAuthorizedUser: errorGetAuthorizedUser,
	}
}

func (service authServiceMock) Authenticate(credentials validators.Credentials, isStaff bool) (*models.User, error) {
	return nil, nil
}

func (service *authServiceMock) Authorize(app *models.App, user *models.User, data validators.OauthAuthorizeData) (string, error) {
	return "", nil
}

func (service *authServiceMock) GenerateTokens(user models.User, scopes []string) services.AuthTokens {
	return services.AuthTokens{AccessToken: "", RefreshToken: ""}
}

func (service authServiceMock) GetAuthorizedUser(accessToken string, scopes []string) (*models.User, error) {
	service.recorder.accessToken = accessToken
	service.recorder.scopes = scopes
	return service.userGetAuthorizedUser, service.errorGetAuthorizedUser
}

func (service authServiceMock) RefreshToken(accessToken string, refreshToken string) (*services.AuthTokens, error) {
	return nil, nil
}

func TestAuthBearerMiddleware(t *testing.T) {
	assert := require.New(t)

	t.Run("Test HasScopes successfully", func(t *testing.T) {
		user := tests.UserFactory()
		authServiceMock := newAuthServiceMock(&user, nil)
		middleware := NewAuthBearerMiddleware(authServiceMock)
		mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())
		token := "mockedtoken"
		mockContext.Request, _ = http.NewRequest("POST", "/", new(bytes.Buffer))
		mockContext.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		scopes := []string{"read:misco"}

		middleware.HasScopes(scopes)(mockContext)
		settedUser := mockContext.MustGet("authorizedUser").(*models.User)

		assert.Equal(authServiceMock.recorder.accessToken, token)
		assert.Equal(authServiceMock.recorder.scopes, scopes)
		assert.Equal(settedUser.Email, user.Email)
	})

	t.Run("Test HasScopes wrong token header", func(t *testing.T) {
		authServiceMock := newAuthServiceMock(nil, nil)
		middleware := NewAuthBearerMiddleware(authServiceMock)
		recorder := httptest.NewRecorder()
		mockContext, _ := gin.CreateTestContext(recorder)
		mockContext.Request, _ = http.NewRequest("POST", "/", new(bytes.Buffer))
		scopes := []string{"read:misco"}

		middleware.HasScopes(scopes)(mockContext)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("Test HasScopes unauthorized user", func(t *testing.T) {
		raisedError := errors.New("wrong")
		authServiceMock := newAuthServiceMock(nil, raisedError)
		middleware := NewAuthBearerMiddleware(authServiceMock)
		recorder := httptest.NewRecorder()
		mockContext, _ := gin.CreateTestContext(recorder)
		token := "mockedtoken"
		mockContext.Request, _ = http.NewRequest("POST", "/", new(bytes.Buffer))
		mockContext.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		scopes := []string{"read:misco"}

		middleware.HasScopes(scopes)(mockContext)

		assert.Equal(recorder.Result().StatusCode, http.StatusForbidden)
	})

	t.Run("Test GetAuthorizedUser successfully", func(t *testing.T) {
		user := tests.UserFactory()
		authServiceMock := newAuthServiceMock(nil, nil)
		middleware := NewAuthBearerMiddleware(authServiceMock)
		mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())
		mockContext.Set("authorizedUser", &user)

		authorizedUser := middleware.GetAuthorizedUser(mockContext)

		assert.Equal(authorizedUser.Email, user.Email)
	})

	t.Run("Test GetAuthorizedUser panics", func(t *testing.T) {
		authServiceMock := newAuthServiceMock(nil, nil)
		middleware := NewAuthBearerMiddleware(authServiceMock)
		mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())

		assert.PanicsWithError(AuthBearerMiddlewareNotCalledError{}.Error(), func() { middleware.GetAuthorizedUser(mockContext) })
	})
}
