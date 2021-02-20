package services

import (
	"errors"
	"gandalf/models"
	"gandalf/validators"
	"os"
	"time"

	set "github.com/deckarep/golang-set"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

/*
Auth scopes
*/
const (
	ScopeUserRead           = "user:read"
	ScopeUserChangePassword = "user:change-password"
	ScopeUserWrite          = "user:write"
	ScopeUserDelete         = "user:delete"
)

/*
accessTokenClaims -> JWT for accessing resources
*/
type accessTokenClaims struct {
	jwt.StandardClaims
	UUID   uuid.UUID
	Email  string
	Scopes set.Set
}

/*
newAccessTokenClaims -> creates claims for the access token from the given
params
*/
func newAccessTokenClaims(user models.User, scopes set.Set, ttl time.Duration) accessTokenClaims {
	return accessTokenClaims{
		UUID:   user.UUID,
		Email:  user.Email,
		Scopes: scopes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl * time.Minute).Unix(),
		},
	}
}

/*
refreshTokenClaims -> JWT for refreshing access token
*/
type refreshTokenClaims struct {
	jwt.StandardClaims
	UUID uuid.UUID
}

/*
newRefreshTokenClaims -> creates claims for the refresh token from the given
params
*/
func newRefreshTokenClaims(user models.User, ttl time.Duration) refreshTokenClaims {
	return refreshTokenClaims{
		UUID: user.UUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl * time.Minute).Unix(),
		},
	}
}

/*
AuthTokens -> contains user tokens for authenticate and refresh
*/
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

/*
IAuthService -> interface for auth service
*/
type IAuthService interface {
	Authenticate(credentials validators.Credentials) (*models.User, error)
	GenerateTokens(user models.User, scopes set.Set) AuthTokens
	GetAuthorizedUser(accessToken string, scopes set.Set) (*models.User, error)
	RefreshToken(accessToken string, refreshToken string) (*AuthTokens, error)
}

/*
AuthService -> auth service
*/
type AuthService struct {
	db        *gorm.DB
	tokenTTL  time.Duration `env:"JWT_TOKEN_TTL"`
	tokenRTTL time.Duration `env:"JWT_TOKEN_RTTL"`
	tokenKey  string        `env:"JWT_TOKEN_KEY"`

	parseTokenWithClaims func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
	newTokenWithClaims   func(method jwt.SigningMethod, claims jwt.Claims) *jwt.Token
	keyfunc              func(token *jwt.Token) (interface{}, error)
}

/*
NewAuthService -> creates a new auth service
*/
func NewAuthService(db *gorm.DB) AuthService {

	keyfunc := func(token *jwt.Token) (interface{}, error) {
		return os.Getenv("JWT_TOKEN_KEY"), nil
	}

	return AuthService{
		db:                   db,
		parseTokenWithClaims: jwt.ParseWithClaims,
		newTokenWithClaims:   jwt.NewWithClaims,
		keyfunc:              keyfunc,
	}
}

/*
getClaims -> Get token claims
*/
func (service AuthService) getClaims(token string, data jwt.Claims, errorOnInvalid bool) error {
	tkn, err := service.parseTokenWithClaims(token, data, service.keyfunc)

	if err != nil || (!tkn.Valid && errorOnInvalid) {
		return AuthorizationError{err}
	}

	return nil
}

/*
signToken -> sign the given token with the private key
*/
func (service AuthService) signToken(token *jwt.Token) string {
	signedToken, err := token.SignedString(service.tokenKey)
	if err != nil {
		panic(err)
	}
	return signedToken
}

/*
Authenticate -> authenticates an user with the given credentials and
returns it
*/
func (service AuthService) Authenticate(credentials validators.Credentials) (*models.User, error) {
	var user models.User
	if err := service.db.Where(&models.User{Email: credentials.Email}).First(&user).Error; err != nil {
		return nil, AuthenticationError{err}
	}

	if !user.VerifyPassword(credentials.Password) {
		return nil, AuthenticationError{nil}
	}

	return &user, nil
}

/*
GenerateTokens -> generate a pair access token for the given user with the
given scopes
*/
func (service AuthService) GenerateTokens(user models.User, scopes set.Set) AuthTokens {
	accessToken := service.signToken(service.newTokenWithClaims(
		jwt.SigningMethodHS256, newAccessTokenClaims(user, scopes, service.tokenTTL),
	))
	refreshToken := service.signToken(service.newTokenWithClaims(
		jwt.SigningMethodHS256, newRefreshTokenClaims(user, service.tokenRTTL),
	))

	return AuthTokens{accessToken, refreshToken}
}

/*
GetAuthorizedUser -> return the user who perform the request if he has
been authorized with the given scopes
*/
func (service AuthService) GetAuthorizedUser(token string, scopes set.Set) (*models.User, error) {
	var accessClaims accessTokenClaims
	err := service.getClaims(token, accessClaims, true)

	if err != nil {
		return nil, err
	}

	if !scopes.IsSubset(accessClaims.Scopes) {
		return nil, AuthorizationError{errors.New("No scopes")}
	}

	var user models.User
	if err := service.db.Where(&models.User{UUID: accessClaims.UUID, Email: accessClaims.Email, Verified: true}).First(&user).Error; err != nil {
		return nil, AuthorizationError{nil}
	}

	user.LastLogin = time.Now()
	service.db.Save(&user)

	return &user, nil
}

/*
RefreshToken -> refresh the access token with his refresh one
*/
func (service AuthService) RefreshToken(accessToken string, refreshToken string) (*AuthTokens, error) {
	var accessClaims accessTokenClaims
	var refreshClaims refreshTokenClaims

	aerr := service.getClaims(accessToken, accessClaims, false)
	rerr := service.getClaims(refreshToken, refreshClaims, true)

	if aerr != nil {
		return nil, AuthenticationError{aerr}
	}
	if rerr != nil {
		return nil, AuthenticationError{rerr}
	}

	if accessClaims.UUID != refreshClaims.UUID {
		return nil, AuthenticationError{errors.New("Access and refresh tokens uuids did not match")}
	}

	accessClaims.ExpiresAt = time.Now().Add(service.tokenTTL * time.Minute).Unix()
	newAccessToken, naerr := service.newTokenWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(service.tokenKey)

	if naerr != nil {
		return nil, AuthenticationError{aerr}
	}

	return &AuthTokens{newAccessToken, refreshToken}, nil
}
