package models

import (
	"gandalf/security"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

const clientSecretLenght = 32

/*
App -> the app itself
*/
type App struct {
	gorm.Model

	// Mandatory fields
	UUID uuid.UUID `gorm:"index:app_uuid;unique;type:uuid;default:uuid_generate_v4()"`

	ClientID     uuid.UUID `gorm:"index:app_client_id;unique;type:uuid;default:uuid_generate_v4()"`
	ClientSecret string    `gorm:"not null"`
	Name         string    `gorm:"not null"`

	// Optional fields
	IconUrl      string
	RedirectUrls pq.StringArray `gorm:"type:text[]"`

	// User
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uint

	// Users that have been sign on the app
	ConnectedUsers []User `gorm:"many2many:user_has_signin_on_app;"`

	// Untracked fields
	secretGenerator security.ISecretGenerator `gorm:"-"`
}

/*
GenerateClientSecret -> generates the client's secret for the
oauth2 connection
*/
func (app *App) generateClientSecret() {
	secret, err := app.secretGenerator.GenerateSecret(clientSecretLenght)
	if err != nil {
		panic(err)
	}
	app.ClientSecret = secret
}

/*
NewApp -> creates a new app
*/
func NewApp(name string, IconUrl string, RedirectUrls []string, user User) App {
	app := App{
		Name:            name,
		IconUrl:         IconUrl,
		RedirectUrls:    RedirectUrls,
		User:            user,
		UserID:          user.ID,
		secretGenerator: security.NewUniformSecret(),
	}
	app.generateClientSecret()
	return app
}
