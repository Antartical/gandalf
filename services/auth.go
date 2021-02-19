package services

import (
	"errors"
	"gandalf/models"
	"gandalf/validators"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

/*
AccessTokenClaims -> JWT for accessing resources
*/
type AccessTokenClaims struct {
	jwt.StandardClaims
	UUID  uuid.UUID
	Email string
}

/*
RefreshTokenClaims -> JWT for refreshing access token
*/
type RefreshTokenClaims struct {
	jwt.StandardClaims
	UUID uuid.UUID
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
	Authenticate(credentials validators.Credentials) (*AuthTokens, error)
	Authorize(accessToken string) (*models.User, error)
	Refresh(accessToken string, refreshToken string) (*AuthTokens, error)
}

/*
AuthService -> auth service
*/
type AuthService struct {
	db        *gorm.DB
	tokenTTL  time.Duration `env:"JWT_TOKEN_TTL"`
	tokenRTTL time.Duration `env:"JWT_TOKEN_RTTL"`
	tokenKey  string        `env:"JWT_TOKEN_KEY"`
}

/*
NewAuthService -> creates a new auth service
*/
func NewAuthService(db *gorm.DB) AuthService {
	return AuthService{db: db}
}

/*
getClaims -> Get token claims
*/
func (service AuthService) getClaims(token string, data jwt.Claims, errorOnInvalid bool) error {
	tkn, err := jwt.ParseWithClaims(token, data, func(token *jwt.Token) (interface{}, error) {
		return service.tokenKey, nil
	})

	if err != nil || (!tkn.Valid && errorOnInvalid) {
		return AuthorizationError{err}
	}

	return nil
}

/*
Authenticate -> authenticate a user and return his identification tokens
*/
func (service AuthService) Authenticate(credentials validators.Credentials) (*AuthTokens, error) {
	var user models.User
	if err := service.db.Where(&models.User{Email: credentials.Email}).First(&user).Error; err != nil {
		return nil, AuthenticationError{err}
	}

	if !user.VerifyPassword(credentials.Password) {
		return nil, AuthenticationError{nil}
	}

	accessClaims := &AccessTokenClaims{
		UUID:  user.UUID,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(service.tokenTTL * time.Minute).Unix(),
		},
	}
	refreshClaims := &RefreshTokenClaims{
		UUID: user.UUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(service.tokenRTTL * time.Minute).Unix(),
		},
	}

	accessToken, aerr := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(service.tokenKey)
	refreshToken, rerr := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(service.tokenKey)

	if aerr != nil {
		return nil, AuthenticationError{aerr}
	}
	if rerr != nil {
		return nil, AuthenticationError{rerr}
	}

	return &AuthTokens{accessToken, refreshToken}, nil
}

/*
Authorize -> authenticates an user by the given token
*/
func (service AuthService) Authorize(token string) (*models.User, error) {
	var accessClaims AccessTokenClaims
	err := service.getClaims(token, accessClaims, true)

	if err != nil {
		return nil, err
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
Refresh -> refresh the access token with his refresh one
*/
func (service AuthService) Refresh(accessToken string, refreshToken string) (*AuthTokens, error) {
	var accessClaims AccessTokenClaims
	var refreshClaims RefreshTokenClaims

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
	newAccessToken, naerr := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(service.tokenKey)

	if naerr != nil {
		return nil, AuthenticationError{aerr}
	}

	return &AuthTokens{newAccessToken, refreshToken}, nil
}
