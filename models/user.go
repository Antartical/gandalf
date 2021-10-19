package models

import (
	"time"

	"gandalf/security"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Represents the basic unit of information for the users that have signed
// into an application by using gandalf
type User struct {
	gorm.Model
	LastLogin time.Time

	// Mandatory fields
	UUID     uuid.UUID `gorm:"index:usr_uuid;unique;type:uuid;default:uuid_generate_v4()"`
	Email    string    `gorm:"not null;index:usr_email;unique"`
	Password string    `gorm:"not null"`
	Name     string    `gorm:"not null"`
	Surname  string    `gorm:"not null"`
	Birthday time.Time `gorm:"not null"`
	Verified bool      `gorm:"default:false"`

	// Optional fields
	Phone string

	// Untracked fields
	hasher security.Hasher `gorm:"-"`

	// Relation fields
	Apps          []App `gorm:"foreignKey:UserID"`
	ConnectedApps []App `gorm:"many2many:user_has_signin_on_app;"`
}

// Set user password by hashing the given one
func (u *User) SetPassword(password string) {
	hash, err := u.hasher.GeneratePassword(password)
	if err != nil {
		panic(err)
	}
	u.Password = string(hash)
}

// Verify if the given password match with the user one
func (u User) VerifyPassword(password string) bool {
	err := u.hasher.VerifyPassword(u.Password, password)
	if err != nil {
		return false
	}
	return true
}

// Gorm hook after find it in the database
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	u.hasher = security.NewBcryptHasher()
	return nil
}

// Creates a new user
func NewUser(email string, password string, name string, surname string, birthday time.Time, phone string) User {
	user := User{
		Email:    email,
		Name:     name,
		Surname:  surname,
		Birthday: birthday,
		Phone:    phone,
		hasher:   security.NewBcryptHasher(),
	}
	user.SetPassword(password)
	return user
}
