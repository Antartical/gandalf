package models

import (
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Claim represents an attempt loggin of the user in the app
// we use it to verify that the given authorization code match
// with the code issure before. A claim can only be used once time
type Claim struct {
	gorm.Model

	// Mandatory fields
	UUID              uuid.UUID      `gorm:"index:app_uuid;unique;type:uuid;default:uuid_generate_v4()"`
	RedirectUrl       string         `gorm:"not null"`
	AuthorizationCode string         `gorm:"not null"`
	Scopes            pq.StringArray `gorm:"type:text[]"`

	// User
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uint

	// App
	App   App `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AppID uint
}

// Creates a new claim
func NewClaim(redirectUrl string, authorizationCode string, scopes []string, user User, app App) Claim {
	return Claim{
		RedirectUrl:       redirectUrl,
		AuthorizationCode: authorizationCode,
		Scopes:            scopes,
		User:              user,
		UserID:            user.ID,
		App:               app,
		AppID:             app.ID,
	}
}
