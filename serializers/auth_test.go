package serializers

import (
	"gandalf/services"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokensSerializer(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		accesToken := "test"
		refreshToken := "test"
		tokens := services.AuthTokens{
			AccessToken:  accesToken,
			RefreshToken: refreshToken,
		}

		serializedTokens := NewTokensSerializer(tokens)

		assert.Equal(serializedTokens.AcessToken, accesToken)
		assert.Equal(serializedTokens.RefreshToken, refreshToken)
	})
}
