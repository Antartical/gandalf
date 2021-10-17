package services

import (
	"gandalf/models"
	"gandalf/validators"

	"gorm.io/gorm"
)

type IAppService interface {
	Create(validators.AppCreateData, models.User) (*models.App, error)
	// Update(uuid.UUID, validators.AppUpdateData) (*models.App, error)
	// Read(uuid uuid.UUID) (*models.App, error)
	// Delete(uuid uuid.UUID) error
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

	return &app, nil
}
