package services

import (
	"gandalf/models"
	"gandalf/validators"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

/*
Claims -> JWT token claims
*/
type Claims struct {
	jwt.StandardClaims
	UUID uuid.UUID
}

/*
IAuthService -> interface for auth service
*/
type IAuthService interface {
	Login(credentials validators.Credentials) (string, error)
	Authenticate(token string) (*models.User, error)
	Refresh(token string) (string, error)
}

/*
AuthService -> auth service
*/
type AuthService struct {
	db       *gorm.DB
	tokenTTL time.Duration `env:"JWT_TOKEN_TTL"`
	tokenKey string        `env:"JWT_TOKEN_KEY"`
}

/*
getClaims -> Get claims from the given token
*/
func (service AuthService) getClaims(token string) (*Claims, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return service.tokenKey, nil
	})

	if err != nil || !tkn.Valid {
		return nil, AuthenticationError{err}
	}

	return claims, nil
}

/*
Login -> authenticate a user and return his identification token
*/
func (service AuthService) Login(credentials validators.Credentials) (string, error) {
	var user models.User
	if err := service.db.Where(&models.User{Email: credentials.Email}).First(&user).Error; err != nil {
		return "", AuthenticationError{err}
	}

	if !user.VerifyPassword(credentials.Password) {
		return "", AuthenticationError{nil}
	}

	expirationTime := time.Now().Add(service.tokenTTL * time.Minute)
	claims := &Claims{
		UUID: user.UUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(service.tokenKey)
}

/*
Authenticate -> authenticates an user by the given token
*/
func (service AuthService) Authenticate(token string) (*models.User, error) {
	claims, err := service.getClaims(token)

	if err != nil {
		return nil, err
	}

	var user models.User
	if err := service.db.Where(&models.User{UUID: claims.UUID, Verified: true}).First(&user).Error; err != nil {
		return nil, AuthenticationError{nil}
	}

	return &user, nil
}
