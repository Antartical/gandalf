package serializers

import (
	"gandalf/services"
)

type TokensSerializer struct {
	AcessToken   string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	RefreshToken string `json:"refresh_token" example:"kpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyf"`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int64  `json:"expires_in" example:"3600"`
}

// Creates a new user serializer
func NewTokensSerializer(tokens services.AuthTokens) TokensSerializer {
	return TokensSerializer{
		AcessToken:   tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(tokens.ExpiresIn),
	}
}
