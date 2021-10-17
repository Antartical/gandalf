package services

import (
	"gandalf/models"
	"gandalf/validators"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type IAppService interface {
	Create(validators.AppCreateData, models.User) (*models.App, error)
	Read(uuid uuid.UUID) (*models.App, error)
	Update(uuid.UUID, validators.AppUpdateData) (*models.App, error)
	Delete(uuid uuid.UUID) error
}

/*
AppService -> app's service
*/
type AppService struct {
	db *gorm.DB
}

/*
NewAppService -> creates a new app service
*/
func NewAppService(db *gorm.DB) AppService {
	return AppService{db}
}

/*
Create -> creates a new app
*/
func (service AppService) Create(appData validators.AppCreateData, user models.User) (*models.App, error) {
	app := models.NewApp(
		appData.Name,
		appData.IconUrl,
		appData.RedirectUrls,
		user,
	)

	if err := service.db.Create(&app).Error; err != nil {
		return nil, AppCreateError{err}
	}

	app.User = user
	return &app, nil
}

/*
Read -> read app from database by his UUID
*/
func (service AppService) Read(uuid uuid.UUID) (*models.App, error) {
	var app models.App
	if err := service.db.Where(&models.App{UUID: uuid}).First(&app).Error; err != nil {
		return nil, AppNotFoundError{err}
	}
	return &app, nil
}

/*
Update -> updates the user which belongs to the given ID according to
the given user data
*/
func (service AppService) Update(uuid uuid.UUID, appData validators.AppUpdateData) (*models.App, error) {
	app, err := service.Read(uuid)
	if err != nil {
		return nil, err
	}

	if appData.Name != "" {
		app.Name = appData.Name
	}

	if appData.IconUrl != "" {
		app.IconUrl = appData.IconUrl
	}

	if len(appData.RedirectUrls) != 0 {
		app.RedirectUrls = appData.RedirectUrls
	}

	service.db.Save(app)
	return app, nil
}

/*
Delete -> set the field `deletes_at` of the app but it will still alive
in database. Soft deleted apps will not appear as result of any query that
not includes `unscoped`
*/
func (service AppService) Delete(uuid uuid.UUID) error {
	if err := service.db.Unscoped().Where(&models.App{UUID: uuid}).Delete(&models.App{}).Error; err != nil {
		return AppNotFoundError{err}
	}
	return nil
}
