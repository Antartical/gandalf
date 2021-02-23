package serializers

import (
	"gandalf/services"
)

type tokenDataSerializer struct {
	AcessToken   string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

/*
TokensSerializer -> token serializer for api output
*/
type TokensSerializer struct {
	ObjectType string              `json:"type"`
	Data       tokenDataSerializer `json:"data"`
}

/*
NewTokensSerializer -> creates a new user serializer and fills it with
the given user data.
*/
func NewTokensSerializer(tokens services.AuthTokens) TokensSerializer {
	return TokensSerializer{
		ObjectType: "tokens",
		Data: tokenDataSerializer{
			AcessToken:   tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}
}
