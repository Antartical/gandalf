package serializers

import (
	"gandalf/helpers"
	"gandalf/models"
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

	t.Run("Test serialize batch", func(t *testing.T) {
		apps := []models.App{tests.AppFactory(), tests.AppFactory(), tests.AppFactory()}
		cursor := helpers.NewCursor(3, 10)
		appSerializer := NewPaginatedAppsSerializer(apps, cursor)
		assert.Equal(len(appSerializer.Data), 3)
	})

	t.Run("Test serialize public", func(t *testing.T) {
		apps := []models.App{tests.AppFactory()}
		cursor := helpers.NewCursor(3, 10)
		appSerializer := NewPaginatedAppsPublicSerializer(apps, cursor)
		assert.Equal(len(appSerializer.Data), 1)
	})
}
