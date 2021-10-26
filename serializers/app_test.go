package serializers

import (
	"gandalf/tests"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppSerializer(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		app := tests.AppFactory()
		appSerializer := NewAppSerializer(app)

		assert.Equal(app.ClientID, appSerializer.Data.ClientID)
		assert.Equal(app.UUID, appSerializer.Data.UUID)
		assert.Equal(app.ClientSecret, appSerializer.Data.ClientSecret)
		assert.Equal(app.Name, appSerializer.Data.Name)
		assert.Equal(app.IconUrl, appSerializer.Data.IconUrl)
		assert.Equal([]string(app.RedirectUrls), appSerializer.Data.RedirectUrls)
	})
}
