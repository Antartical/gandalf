package models

import (
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
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
		IconUrl := "http://fakeicon.ico"
		RedirectUrls := []string{"FakeUri"}
		user := User{}
		user.ID = uint(faker.Number().NumberInt(3))

		app := NewApp(name, IconUrl, RedirectUrls, user)

		assert.Equal(app.Name, name)
		assert.Equal(app.IconUrl, IconUrl)
		assert.Equal(app.RedirectUrls, (pq.StringArray)(RedirectUrls))
		assert.Equal(app.UserID, user.ID)
	})

	t.Run("Test constructor fail", func(t *testing.T) {
		expectedError := errors.New("Whoops")
		name := "Fake app"
		IconUrl := "http://fakeicon.ico"
		RedirectUrls := []string{"FakeUri"}
		user := User{}
		user.ID = uint(faker.Number().NumberInt(3))
		secretGenerator := mockedSecretGenerator{
			generateSecretError: expectedError,
		}

		app := App{
			Name:            name,
			IconUrl:         IconUrl,
			RedirectUrls:    RedirectUrls,
			UserID:          user.ID,
			secretGenerator: &secretGenerator,
		}

		assert.PanicsWithError(expectedError.Error(), func() { app.generateClientSecret() })
	})
}
