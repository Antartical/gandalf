package services

import (
	"gandalf/helpers"
	"gandalf/tests"
	"gandalf/validators"
	"testing"

	"github.com/gofrs/uuid"
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
		user := tests.UserFactory()
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
		user := tests.UserFactory()
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

func TestAppServiceRead(t *testing.T) {
	assert := require.New(t)

	t.Run("Test read app successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := AppService{db}

		app := tests.AppFactory()
		db.Create(&app)

		readApp, err := service.Read(app.UUID)

		assert.NoError(err)
		assert.Equal(app.ID, readApp.ID)

		db.Unscoped().Delete(&app.User)
	})

	t.Run("Test read app not found error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := AppService{db}
		uuid, _ := uuid.NewV4()

		_, err := service.Read(uuid)
		assert.Error(err, AppNotFoundError{nil}.Error())
	})

}

func TestAppServiceReadByClientID(t *testing.T) {
	assert := require.New(t)

	t.Run("Test read app by client successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := AppService{db}

		app := tests.AppFactory()
		db.Create(&app)

		readApp, err := service.ReadByClientID(app.ClientID)

		assert.NoError(err)
		assert.Equal(app.ID, readApp.ID)

		db.Unscoped().Delete(&app.User)
	})

	t.Run("Test read app by client id not found error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := AppService{db}
		clientID, _ := uuid.NewV4()

		_, err := service.ReadByClientID(clientID)
		assert.Error(err, AppNotFoundError{nil}.Error())
	})

}

func TestAppServiceUpdate(t *testing.T) {
	assert := require.New(t)

	t.Run("Test update app successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := AppService{db}

		app := tests.AppFactory()
		db.Create(&app)

		name := faker.Company().Name()
		iconUrl := faker.Internet().Url()
		redirectUrls := []string{faker.Internet().Url()}

		appData := validators.AppUpdateData{
			Name:         name,
			IconUrl:      iconUrl,
			RedirectUrls: redirectUrls,
		}

		updatedApp, err := service.Update(app.UUID, appData)

		assert.NoError(err)
		assert.Equal(name, updatedApp.Name)
		assert.Equal(iconUrl, updatedApp.IconUrl)
		assert.Equal(redirectUrls[0], updatedApp.RedirectUrls[0])

		db.Unscoped().Delete(&app.User)
	})

	t.Run("Test update app not found error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := AppService{db}

		uuid, _ := uuid.NewV4()
		name := faker.Company().Name()
		iconUrl := faker.Internet().Url()
		redirectUrls := []string{faker.Internet().Url()}
		appData := validators.AppUpdateData{
			Name:         name,
			IconUrl:      iconUrl,
			RedirectUrls: redirectUrls,
		}

		_, err := service.Update(uuid, appData)

		assert.Error(err, AppNotFoundError{nil}.Error())
	})

}

func TestAppServiceDelete(t *testing.T) {
	assert := require.New(t)

	t.Run("Test delete app successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := AppService{db}

		app := tests.AppFactory()
		db.Create(&app)

		err := service.Delete(app.UUID)
		assert.NoError(err)
	})

	t.Run("Test delete app error not found", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := AppService{db}

		app := tests.AppFactory()

		err := service.Delete(app.UUID)
		assert.Error(err, AppNotFoundError{nil}.Error())
	})

}

func TestAppServiceListApps(t *testing.T) {
	assert := require.New(t)

	t.Run("Test list apps", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := AppService{db}

		app := tests.AppFactory()
		db.Create(&app)

		cursor := helpers.NewCursor(0, 30)
		assert.Equal(app.ID, service.ListApps(app.User, &cursor)[0].ID)
		db.Delete(&app)
	})

}

func TestAppServiceListConnectedApps(t *testing.T) {
	assert := require.New(t)

	t.Run("Test list apps", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := AppService{db}
		app := tests.AppFactory()
		db.Create(&app)
		db.Model(&app).Association("ConnectedUsers").Append(&app.User)

		cursor := helpers.NewCursor(0, 30)
		assert.Equal(app.ID, service.ListConnectedApps(app.User, &cursor)[0].ID)
		db.Delete(&app)
	})

}
