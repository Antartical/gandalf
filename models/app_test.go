package models

import (
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type mockedSecretGenerator struct {
	generateSecretError error
}

func (secretGenerator *mockedSecretGenerator) GenerateSecret(lenght int) (string, error) {
	return "fakeSecret", secretGenerator.generateSecretError
}

func TestAppModel(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor success", func(t *testing.T) {
		name := "Fake app"
		iconUri := "http://fakeicon.ico"
		redirectUris := []string{"FakeUri"}
		user := UserFactory()

		app := NewApp(name, iconUri, redirectUris, user)

		assert.Equal(app.Name, name)
		assert.Equal(app.IconUri, iconUri)
		assert.Equal(app.RedirectUris, (pq.StringArray)(redirectUris))
		assert.Equal(app.User.ID, user.ID)
	})

	t.Run("Test constructor fail", func(t *testing.T) {
		expectedError := errors.New("Whoops")
		name := "Fake app"
		iconUri := "http://fakeicon.ico"
		redirectUris := []string{"FakeUri"}
		user := UserFactory()
		secretGenerator := mockedSecretGenerator{
			generateSecretError: expectedError,
		}

		app := App{
			Name:            name,
			IconUri:         iconUri,
			RedirectUris:    redirectUris,
			User:            user,
			UserID:          user.ID,
			secretGenerator: &secretGenerator,
		}

		assert.PanicsWithError(expectedError.Error(), func() { app.generateClientSecret() })
	})
}
