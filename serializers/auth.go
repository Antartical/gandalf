package serializers

import (
	"gandalf/services"
)

type tokenDataSerializer struct {
	AcessToken   string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	RefreshToken string `json:"refresh_token" example:"kpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyf"`
}

// Bearer auth token serialization struct
type TokensSerializer struct {
	ObjectType string              `json:"type" example:"tokens"`
	Data       tokenDataSerializer `json:"data"`
}

// Creates a new user serializer
func NewTokensSerializer(tokens services.AuthTokens) TokensSerializer {
	return TokensSerializer{
		ObjectType: "tokens",
		Data: tokenDataSerializer{
			AcessToken:   tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}
}
