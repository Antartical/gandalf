package services

import (
	"gandalf/models"
	"gandalf/tests"
	"gandalf/validators"
	"testing"

	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestAppServiceConstructor(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		service := NewAppService(db)

		assert.Equal(service.db, db)
	})
}

func TestAppServiceCreate(t *testing.T) {
	assert := require.New(t)

	t.Run("Test app create successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := AppService{db}
		user := models.UserFactory()
		db.Create(&user)

		name := faker.Company().Name()
		iconUrl := faker.Internet().Url()
		redirectUrls := []string{faker.Internet().Url()}

		appData := validators.AppCreateData{
			Name:         name,
			IconUrl:      iconUrl,
			RedirectUrls: redirectUrls,
		}

		app, err := service.Create(appData, user)

		assert.NoError(err)
		assert.Equal(name, app.Name)
		assert.Equal(iconUrl, app.IconUrl)
		assert.Equal(redirectUrls[0], app.RedirectUrls[0])

		db.Unscoped().Delete(&app)
		db.Unscoped().Delete(&user)
	})

	t.Run("Test app create database error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := AppService{db}
		user := models.UserFactory()
		user.ID = uint(faker.Number().NumberInt(3))
		name := faker.Company().Name()
		iconUrl := faker.Internet().Url()
		redirectUrls := []string{faker.Internet().Url()}

		appData := validators.AppCreateData{
			Name:         name,
			IconUrl:      iconUrl,
			RedirectUrls: redirectUrls,
		}

		_, err := service.Create(appData, user)
		assert.Error(err, AppCreateError{nil}.Error())
	})

}
