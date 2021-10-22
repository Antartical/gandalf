package services

import (
	"errors"
	"gandalf/bindings"
	"gandalf/helpers"
	"gandalf/models"
	"gandalf/security"
	"gandalf/validators"
	"os"
	"strconv"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// JWT for accessing resources
type accessTokenClaims struct {
	jwt.StandardClaims
	UUID   uuid.UUID
	Email  string
	Scopes []string
}

// Creates claims for the access token from the given params
func newAccessTokenClaims(user models.User, scopes []string, ttl time.Duration) accessTokenClaims {
	return accessTokenClaims{
		UUID:   user.UUID,
		Email:  user.Email,
		Scopes: scopes,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl * time.Minute).Unix(),
		},
	}
}

// JWT for refreshing access token
type refreshTokenClaims struct {
	jwt.StandardClaims
	UUID uuid.UUID
}

// Creates claims for the refresh token from the given params
func newRefreshTokenClaims(user models.User, ttl time.Duration) refreshTokenClaims {
	return refreshTokenClaims{
		UUID: user.UUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl * time.Minute).Unix(),
		},
	}
}

// Contains user tokens for authenticate and refresh
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
}

// Interface for auth service
type IAuthService interface {
	Authenticate(credentials validators.Credentials, isStaff bool) (*models.User, error)
	GenerateTokens(user models.User, scopes []string) AuthTokens
	GetAuthorizedUser(accessToken string, scopes []string) (*models.User, error)
	RefreshToken(accessToken string, refreshToken string) (*AuthTokens, error)
	Authorize(*models.App, *models.User, validators.OauthAuthorizeData) (string, error)
	ExchangeOauthToken(models.App, validators.OauthExchangeToken) (*AuthTokens, error)
}

// Auth service
type AuthService struct {
	db        *gorm.DB
	tokenTTL  time.Duration `env:"JWT_TOKEN_TTL"`
	tokenRTTL time.Duration `env:"JWT_TOKEN_RTTL"`
	tokenKey  interface{}   `env:"JWT_TOKEN_KEY"`

	parseTokenWithClaims func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error)
	newTokenWithClaims   func(method jwt.SigningMethod, claims jwt.Claims) *jwt.Token
	keyfunc              func(token *jwt.Token) (interface{}, error)
}

// Creates a new auth service
func NewAuthService(db *gorm.DB) AuthService {

	keyfunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_TOKEN_KEY")), nil
	}

	tokenTTL, _ := strconv.Atoi(os.Getenv("JWT_TOKEN_TTL"))
	tokenRTTL, _ := strconv.Atoi(os.Getenv("JWT_TOKEN_RTTL"))

	return AuthService{
		db:                   db,
		tokenTTL:             time.Duration(tokenTTL),
		tokenRTTL:            time.Duration(tokenRTTL),
		tokenKey:             []byte(os.Getenv("JWT_TOKEN_KEY")),
		parseTokenWithClaims: jwt.ParseWithClaims,
		newTokenWithClaims:   jwt.NewWithClaims,
		keyfunc:              keyfunc,
	}
}

// Get token claims
func (service AuthService) getClaims(token string, data jwt.Claims, errorOnInvalid bool) error {
	tkn, err := service.parseTokenWithClaims(token, data, service.keyfunc)

	if err != nil || (!tkn.Valid && errorOnInvalid) {
		return AuthorizationError{err}
	}

	return nil
}

//Sign the given token with the private key
func (service AuthService) signToken(token *jwt.Token) string {
	signedToken, err := token.SignedString(service.tokenKey)
	if err != nil {
		panic(err)
	}
	return signedToken
}

// Authenticates an user with the given credentials and returns it
func (service AuthService) Authenticate(credentials validators.Credentials, isStaff bool) (*models.User, error) {
	var user models.User
	if err := service.db.Where(&models.User{Email: credentials.Email, Verified: true, Staff: isStaff}).First(&user).Error; err != nil {
		return nil, AuthenticationError{err}
	}

	if !user.VerifyPassword(credentials.Password) {
		return nil, AuthenticationError{nil}
	}

	return &user, nil
}

// Generate a pair access token for the given user with the given scopes
func (service AuthService) GenerateTokens(user models.User, scopes []string) AuthTokens {
	accessToken := service.signToken(service.newTokenWithClaims(
		jwt.SigningMethodHS256, newAccessTokenClaims(user, scopes, service.tokenTTL),
	))
	refreshToken := service.signToken(service.newTokenWithClaims(
		jwt.SigningMethodHS256, newRefreshTokenClaims(user, service.tokenRTTL),
	))

	return AuthTokens{accessToken, refreshToken, service.tokenTTL}
}

// Return the user who perform the request if he has
// been authorized with the given scopes
func (service AuthService) GetAuthorizedUser(token string, scopes []string) (*models.User, error) {
	accessClaims := &accessTokenClaims{}
	err := service.getClaims(token, accessClaims, true)
	verified := true

	if err != nil {
		return nil, err
	}

	mandatoryScopes := mapset.NewSet()
	for _, elem := range scopes {
		mandatoryScopes.Add(elem)
	}
	tokenScopes := mapset.NewSet()
	for _, elem := range accessClaims.Scopes {
		tokenScopes.Add(elem)
	}

	// It's mandatory to search on verified users, except on the verification
	// endpoint
	if mandatoryScopes.Contains(security.ScopeUserVerify) {
		verified = false
	}

	if !mandatoryScopes.IsSubset(tokenScopes) {
		return nil, AuthorizationError{errors.New("Unauthorized")}
	}

	var user models.User
	if err := service.db.Where(&models.User{UUID: accessClaims.UUID, Email: accessClaims.Email, Verified: verified}).First(&user).Error; err != nil {
		return nil, AuthorizationError{errors.New("Related user does not exist")}
	}

	user.LastLogin = time.Now()
	service.db.Save(&user)

	return &user, nil
}

// Refresh the access token with his refresh one
func (service AuthService) RefreshToken(accessToken string, refreshToken string) (*AuthTokens, error) {
	accessClaims := &accessTokenClaims{}
	refreshClaims := &refreshTokenClaims{}

	aerr := service.getClaims(accessToken, accessClaims, false)
	rerr := service.getClaims(refreshToken, refreshClaims, true)

	if aerr != nil || rerr != nil {
		return nil, AuthenticationError{errors.New("Unrecognized token")}
	}

	if accessClaims.UUID != refreshClaims.UUID {
		return nil, AuthenticationError{errors.New("Unrelated access and refresh token")}
	}

	accessClaims.ExpiresAt = time.Now().Add(service.tokenTTL * time.Minute).Unix()
	newAccessToken := service.signToken(service.newTokenWithClaims(jwt.SigningMethodHS256, accessClaims))

	return &AuthTokens{newAccessToken, refreshToken, service.tokenTTL}, nil
}

// Associate the given app with the given app in order to save that the user
// has signin on the given app. Returns the authorization code and error.
func (service AuthService) Authorize(app *models.App, user *models.User, data validators.OauthAuthorizeData) (string, error) {

	if !helpers.PqStringArrayContains(app.RedirectUrls, data.RedirectURI) {
		return "", RedirectUriDoesNotMatch{redirectUri: data.RedirectURI}
	}

	authorizationCode := service.GenerateTokens(*user, []string{security.ScopeUserAuthorizationCode}).AccessToken
	claim := models.NewClaim(
		data.RedirectURI,
		authorizationCode,
		bindings.ScopeArrayToStringArray(data.Scopes),
		*user,
		*app,
	)

	service.db.Create(&claim)
	service.db.Model(app).Association("ConnectedUsers").Append(user)
	service.db.Model(user).Association("ConnectedApps").Find(&user.ConnectedApps)
	service.db.Model(app).Association("ConnectedUser").Find(&app.ConnectedUsers)

	return authorizationCode, nil
}

// Produces an access token with the requested scopes if the given data belongs to the
// created claim. Otherwise an error will be returned
func (service AuthService) ExchangeOauthToken(app models.App, data validators.OauthExchangeToken) (*AuthTokens, error) {
	if !helpers.PqStringArrayContains(app.RedirectUrls, data.RedirectUrl) {
		return nil, RedirectUriDoesNotMatch{redirectUri: data.RedirectUrl}
	}

	user, err := service.GetAuthorizedUser(data.AuthorizationCode, []string{security.ScopeUserAuthorizationCode})
	if err != nil {
		return nil, err
	}

	clause := &models.Claim{
		AuthorizationCode: data.AuthorizationCode,
		RedirectUrl:       data.RedirectUrl,
		AppID:             app.ID,
		UserID:            user.ID,
	}

	var claim models.Claim
	if err := service.db.Where(clause).First(&claim).Error; err != nil {
		return nil, ClaimDoesNotExist{err}
	}

	tokens := service.GenerateTokens(*user, claim.Scopes)
	return &tokens, nil
}
