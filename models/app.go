package models

import (
	"gandalf/security"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

const clientSecretLengt = 32

/*
App -> the app itself
*/
type App struct {
	gorm.Model

	// Mandatory fields
	UUID uuid.UUID `gorm:"index:app_uuid;unique;type:uuid;default:uuid_generate_v4()"`

	ClientID     string `gorm:"index:app_uuid;unique;type:text:default:uuid_generate_v4()"`
	ClientSecret string `gorm:"index:usr_uuid;type:text"`
	Name         string `gorm:"not null"`

	// Optional fields
	IconUri      string
	RedirectUris pq.StringArray `gorm:"type:text[]"`

	// User
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uint
}

/*
NewApp -> creates a new app
*/
func NewApp(name string, iconUri string, redirectUris []string, user User, secretGenerator security.ISecretGenerator) (*App, error) {
	secret, err := secretGenerator.GenerateSecret(clientSecretLengt)
	if err != nil {
		return nil, err
	}

	app := App{
		ClientSecret: secret,
		Name:         name,
		IconUri:      iconUri,
		RedirectUris: redirectUris,
		User:         user,
		UserID:       user.ID,
	}

	return &app, nil
}
