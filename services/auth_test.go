package services

import (
	"gandalf/bindings"
	"gandalf/models"
	"gandalf/security"
	"gandalf/tests"
	"gandalf/validators"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"syreclabs.com/go/faker"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

type parseClaimsRecorder struct {
	tokenString string
	claims      jwt.Claims
}

func TestAuthService(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)

		expectedTokenTTL, _ := strconv.Atoi(os.Getenv("JWT_TOKEN_TTL"))
		expectedTokenRTTL, _ := strconv.Atoi(os.Getenv("JWT_TOKEN_RTTL"))

		assert.Equal(authService.tokenKey, []byte(os.Getenv("JWT_TOKEN_KEY")))
		assert.Equal(authService.tokenTTL, time.Duration(expectedTokenTTL))
		assert.Equal(authService.tokenRTTL, time.Duration(expectedTokenRTTL))
	})

	t.Run("Test getClaims successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)

		user := tests.UserFactory()
		user.UUID, _ = uuid.NewV4()
		scopes := []string{security.ScopeUserRead}
		mockToken := authService.signToken(authService.newTokenWithClaims(
			jwt.SigningMethodHS256, newAccessTokenClaims(user, scopes, authService.tokenTTL),
		))

		claims := &accessTokenClaims{}
		err := authService.getClaims(mockToken, claims, true)

		assert.NoError(err)
		assert.Equal(claims.UUID, user.UUID)
		assert.Equal(claims.Scopes, scopes)
	})

	t.Run("Test getClaims error", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)
		raisedError := errors.New("wrong")
		recorder := &parseClaimsRecorder{}
		authService.parseTokenWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
			recorder.tokenString = tokenString
			recorder.claims = claims
			return nil, raisedError
		}

		user := tests.UserFactory()
		user.UUID, _ = uuid.NewV4()
		scopes := []string{security.ScopeUserRead}
		mockToken := authService.signToken(authService.newTokenWithClaims(
			jwt.SigningMethodHS256, newAccessTokenClaims(user, scopes, authService.tokenTTL),
		))

		claims := &accessTokenClaims{}
		err := authService.getClaims(mockToken, claims, true)

		assert.Error(err, AuthorizationError{raisedError}.Error())
		assert.Equal(recorder.tokenString, mockToken)
		assert.Equal(recorder.claims, claims)
	})

	t.Run("Test signToken successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)

		user := tests.UserFactory()
		user.UUID, _ = uuid.NewV4()
		scopes := []string{security.ScopeUserRead}

		assert.NotPanics(func() {
			authService.signToken(authService.newTokenWithClaims(
				jwt.SigningMethodHS256, newAccessTokenClaims(user, scopes, authService.tokenTTL),
			))
		})
	})

	t.Run("Test signToken panics", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)
		authService.tokenKey = "string"

		user := tests.UserFactory()
		user.UUID, _ = uuid.NewV4()
		scopes := []string{security.ScopeUserRead}

		assert.Panics(func() {
			authService.signToken(authService.newTokenWithClaims(
				jwt.SigningMethodHS256, newAccessTokenClaims(user, scopes, authService.tokenTTL),
			))
		})
	})

	t.Run("Test Authenticate successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		authService := NewAuthService(db)
		plainPassword := "testestestestest"
		user := tests.UserFactory()
		user.SetPassword(plainPassword)
		user.Verified = true
		db.Create(&user)

		credentials := validators.Credentials{
			Email:    user.Email,
			Password: plainPassword,
		}

		authenticatedUser, err := authService.Authenticate(credentials, false)

		assert.NoError(err)
		assert.Equal(authenticatedUser.UUID, user.UUID)
		db.Unscoped().Delete(&user)
	})

	t.Run("Test Authenticate error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		authService := NewAuthService(db)
		plainPassword := "testestestestest"
		user := tests.UserFactory()

		credentials := validators.Credentials{
			Email:    user.Email,
			Password: plainPassword,
		}

		_, err := authService.Authenticate(credentials, false)

		assert.Error(err, AuthenticationError{nil}.Error())
	})

	t.Run("Test Authenticate error verify", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		authService := NewAuthService(db)
		plainInventedPassword := "okkookkookkookkookko"
		user := tests.UserFactory()
		user.Verified = true
		db.Create(&user)

		credentials := validators.Credentials{
			Email:    user.Email,
			Password: plainInventedPassword,
		}

		_, err := authService.Authenticate(credentials, false)

		assert.Error(err, AuthenticationError{nil}.Error())
		db.Unscoped().Delete(&user)
	})

	t.Run("Test Generate tokens", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)

		user := tests.UserFactory()
		user.UUID, _ = uuid.NewV4()

		assert.NotNil(authService.GenerateTokens(user, []string{}))
	})

	t.Run("Test GetAuthorizedUser successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		authService := NewAuthService(db)
		scopes := []string{security.ScopeUserRead, security.ScopeUserVerify}
		user := tests.UserFactory()
		user.Verified = true
		db.Create(&user)

		tokens := authService.GenerateTokens(user, scopes)
		authorizedUser, err := authService.GetAuthorizedUser(tokens.AccessToken, scopes)

		assert.NoError(err)
		assert.Equal(authorizedUser.UUID, user.UUID)

		db.Unscoped().Delete(&user)
	})

	t.Run("Test GetAuthorizedUser getClaims error", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)
		raisedError := errors.New("wrong")
		authService.parseTokenWithClaims = func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
			return nil, raisedError
		}
		scopes := []string{security.ScopeUserRead}
		user := tests.UserFactory()

		tokens := authService.GenerateTokens(user, scopes)
		_, err := authService.GetAuthorizedUser(tokens.AccessToken, scopes)

		assert.Error(err, AuthenticationError{raisedError}.Error())
	})

	t.Run("Test GetAuthorizedUser no scopes error", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)
		scopes := []string{security.ScopeUserRead}
		otherScopes := []string{security.ScopeUserWrite}
		user := tests.UserFactory()

		tokens := authService.GenerateTokens(user, scopes)
		_, err := authService.GetAuthorizedUser(tokens.AccessToken, otherScopes)

		assert.Error(err, AuthenticationError{nil}.Error())
	})

	t.Run("Test GetAuthorizedUser no user error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		authService := NewAuthService(db)
		scopes := []string{security.ScopeUserRead}
		user := tests.UserFactory()

		tokens := authService.GenerateTokens(user, scopes)
		_, err := authService.GetAuthorizedUser(tokens.AccessToken, scopes)

		assert.Error(err, AuthenticationError{nil}.Error())
	})

	t.Run("Test RefreshToken successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)
		scopes := []string{security.ScopeUserRead}
		user := tests.UserFactory()

		tokens := authService.GenerateTokens(user, scopes)
		newTokens, err := authService.RefreshToken(tokens.AccessToken, tokens.RefreshToken)

		assert.NoError(err)
		assert.Equal(newTokens.RefreshToken, tokens.RefreshToken)
	})

	t.Run("Test RefreshToken unrecognized token", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)

		_, err := authService.RefreshToken("false", "misco")

		assert.Error(err, AuthenticationError{nil}.Error())
	})

	t.Run("Test RefreshToken unrelates tokens", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)
		scopes := []string{security.ScopeUserRead}
		user := tests.UserFactory()

		accessToken := authService.GenerateTokens(user, scopes).AccessToken
		user.UUID, _ = uuid.NewV4()
		refreshToken := authService.GenerateTokens(user, scopes).RefreshToken

		_, err := authService.RefreshToken(accessToken, refreshToken)

		assert.Error(err, AuthenticationError{}.Error())
	})

}

func TestAppServiceAuthorize(t *testing.T) {
	assert := require.New(t)

	t.Run("Test authorize success", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := NewAuthService(db)

		app := tests.AppFactory()
		user := tests.UserFactory()
		db.Create(&app)
		db.Create(&user)

		input := validators.OauthAuthorizeData{
			ClientID:    app.ClientID.String(),
			RedirectURI: app.RedirectUrls[0],
			Scopes:      []bindings.Scope{security.ScopeUserRead},
			State:       "state",
		}

		service.Authorize(&app, &user, input)
		assert.Equal(1, len(user.ConnectedApps))
		assert.Equal(1, len(app.ConnectedUsers))
		assert.Equal(user.ID, app.ConnectedUsers[0].ID)
		assert.Equal(app.ID, user.ConnectedApps[0].ID)

		db.Delete(&app)
		db.Delete(&user)
	})

	t.Run("Test authorize fail", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := NewAuthService(db)
		fakeRedirectUri := faker.Internet().Url()
		expectedError := RedirectUriDoesNotMatch{
			redirectUri: fakeRedirectUri,
		}
		//expectedMsg := fmt.Sprintf("Redirect uri is not registered for the app, %s", fakeRedirectUri)

		app := tests.AppFactory()
		user := tests.UserFactory()
		db.Create(&app)
		db.Create(&user)

		input := validators.OauthAuthorizeData{
			ClientID:    app.ClientID.String(),
			RedirectURI: fakeRedirectUri,
			Scopes:      []bindings.Scope{security.ScopeUserRead},
			State:       "state",
		}

		_, err := service.Authorize(&app, &user, input)
		assert.Error(err, expectedError, expectedError.Error())

		db.Delete(&app)
		db.Delete(&user)
	})
}

func TestAppServiceExchangeOauthToken(t *testing.T) {
	assert := require.New(t)

	t.Run("Test Exchange token success", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := NewAuthService(db)
		app := tests.AppFactory()
		user := tests.UserFactory()
		user.Verified = true
		scopes := []string{security.ScopeUserRead}
		db.Create(&app)
		db.Create(&user)

		tokens := service.GenerateTokens(user, []string{security.ScopeUserAuthorizationCode})

		data := validators.OauthExchangeToken{
			GrantType:         "authorization_code",
			ClientID:          app.ClientID.String(),
			ClientSecret:      app.ClientSecret,
			AuthorizationCode: tokens.AccessToken,
			RedirectUrl:       app.RedirectUrls[0],
		}
		claim := models.NewClaim(
			data.RedirectUrl,
			data.AuthorizationCode,
			scopes,
			user,
			app,
		)
		db.Create(&claim)

		resultTokens, err := service.ExchangeOauthToken(app, data)

		assert.Nil(err)
		assert.NotNil(resultTokens)

		db.Delete(&claim)
		db.Delete(&app)
		db.Delete(&user)
	})

	t.Run("Test Exchange cannot get user", func(t *testing.T) {
		expectedError := AuthorizationError{}
		db := tests.NewTestDatabase(false)
		service := NewAuthService(db)
		app := tests.AppFactory()
		user := tests.UserFactory()
		user.Verified = true
		scopes := []string{security.ScopeUserRead}
		db.Create(&app)
		db.Create(&user)

		tokens := service.GenerateTokens(user, scopes)

		data := validators.OauthExchangeToken{
			GrantType:         "authorization_code",
			ClientID:          app.ClientID.String(),
			ClientSecret:      app.ClientSecret,
			AuthorizationCode: tokens.AccessToken,
			RedirectUrl:       app.RedirectUrls[0],
		}
		claim := models.NewClaim(
			data.RedirectUrl,
			data.AuthorizationCode,
			scopes,
			user,
			app,
		)
		db.Create(&claim)

		_, err := service.ExchangeOauthToken(app, data)

		assert.Error(expectedError, err)

		db.Delete(&claim)
		db.Delete(&app)
		db.Delete(&user)
	})

	t.Run("Test Exchange token redirect url does not match", func(t *testing.T) {
		expectedError := RedirectUriDoesNotMatch{}
		db := tests.NewTestDatabase(false)
		service := NewAuthService(db)
		app := tests.AppFactory()
		user := tests.UserFactory()
		user.Verified = true
		scopes := []string{security.ScopeUserRead}
		db.Create(&app)
		db.Create(&user)

		tokens := service.GenerateTokens(user, scopes)

		data := validators.OauthExchangeToken{
			GrantType:         "authorization_code",
			ClientID:          app.ClientID.String(),
			ClientSecret:      app.ClientSecret,
			AuthorizationCode: tokens.AccessToken,
			RedirectUrl:       faker.Internet().Url(),
		}
		claim := models.NewClaim(
			data.RedirectUrl,
			data.AuthorizationCode,
			scopes,
			user,
			app,
		)
		db.Create(&claim)

		_, err := service.ExchangeOauthToken(app, data)

		assert.Error(expectedError, err)

		db.Delete(&claim)
		db.Delete(&app)
		db.Delete(&user)
	})

	t.Run("Test Exchange token claim does not exist", func(t *testing.T) {
		expectedError := ClaimDoesNotExist{}
		db := tests.NewTestDatabase(false)
		service := NewAuthService(db)
		app := tests.AppFactory()
		user := tests.UserFactory()
		user.Verified = true
		db.Create(&app)
		db.Create(&user)

		tokens := service.GenerateTokens(user, []string{security.ScopeUserAuthorizationCode})

		data := validators.OauthExchangeToken{
			GrantType:         "authorization_code",
			ClientID:          app.ClientID.String(),
			ClientSecret:      app.ClientSecret,
			AuthorizationCode: tokens.AccessToken,
			RedirectUrl:       app.RedirectUrls[0],
		}

		_, err := service.ExchangeOauthToken(app, data)

		assert.Error(err, expectedError.Error())

		db.Delete(&app)
		db.Delete(&user)
	})
}
