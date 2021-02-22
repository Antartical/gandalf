package services

import (
	"gandalf/tests"
	"gandalf/validators"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

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

		user := userFactory()
		user.UUID, _ = uuid.NewV4()
		scopes := []string{ScopeUserRead}
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

		user := userFactory()
		user.UUID, _ = uuid.NewV4()
		scopes := []string{ScopeUserRead}
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

		user := userFactory()
		user.UUID, _ = uuid.NewV4()
		scopes := []string{ScopeUserRead}

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

		user := userFactory()
		user.UUID, _ = uuid.NewV4()
		scopes := []string{ScopeUserRead}

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
		user := userFactory()
		user.Verified = true
		db.Create(&user)

		credentials := validators.Credentials{
			Email:    user.Email,
			Password: plainPassword,
		}

		authenticatedUser, err := authService.Authenticate(credentials)

		assert.NoError(err)
		assert.Equal(authenticatedUser.UUID, user.UUID)
		db.Unscoped().Delete(&user)
	})

	t.Run("Test Authenticate error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		authService := NewAuthService(db)
		plainPassword := "testestestestest"
		user := userFactory()

		credentials := validators.Credentials{
			Email:    user.Email,
			Password: plainPassword,
		}

		_, err := authService.Authenticate(credentials)

		assert.Error(err, AuthenticationError{nil}.Error())
	})

	t.Run("Test Authenticate error verify", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		authService := NewAuthService(db)
		plainInventedPassword := "okkookkookkookkookko"
		user := userFactory()
		user.Verified = true
		db.Create(&user)

		credentials := validators.Credentials{
			Email:    user.Email,
			Password: plainInventedPassword,
		}

		_, err := authService.Authenticate(credentials)

		assert.Error(err, AuthenticationError{nil}.Error())
		db.Unscoped().Delete(&user)
	})

	t.Run("Test Generate tokens", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)

		user := userFactory()
		user.UUID, _ = uuid.NewV4()

		assert.NotNil(authService.GenerateTokens(user, []string{}))
	})

	t.Run("Test GetAuthorizedUser successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		authService := NewAuthService(db)
		scopes := []string{ScopeUserRead}
		user := userFactory()
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
		scopes := []string{ScopeUserRead}
		user := userFactory()

		tokens := authService.GenerateTokens(user, scopes)
		_, err := authService.GetAuthorizedUser(tokens.AccessToken, scopes)

		assert.Error(err, AuthenticationError{raisedError}.Error())
	})

	t.Run("Test GetAuthorizedUser no scopes error", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)
		scopes := []string{ScopeUserRead}
		otherScopes := []string{ScopeUserWrite}
		user := userFactory()

		tokens := authService.GenerateTokens(user, scopes)
		_, err := authService.GetAuthorizedUser(tokens.AccessToken, otherScopes)

		assert.Error(err, AuthenticationError{nil}.Error())
	})

	t.Run("Test GetAuthorizedUser no user error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		authService := NewAuthService(db)
		scopes := []string{ScopeUserRead}
		user := userFactory()

		tokens := authService.GenerateTokens(user, scopes)
		_, err := authService.GetAuthorizedUser(tokens.AccessToken, scopes)

		assert.Error(err, AuthenticationError{nil}.Error())
	})

	t.Run("Test RefreshToken successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		authService := NewAuthService(db)
		scopes := []string{ScopeUserRead}
		user := userFactory()

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
		scopes := []string{ScopeUserRead}
		user := userFactory()

		accessToken := authService.GenerateTokens(user, scopes).AccessToken
		user.UUID, _ = uuid.NewV4()
		refreshToken := authService.GenerateTokens(user, scopes).RefreshToken

		_, err := authService.RefreshToken(accessToken, refreshToken)

		assert.Error(err, AuthenticationError{}.Error())
	})

}
