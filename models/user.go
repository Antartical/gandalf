package models

import (
	"time"

	"gandalf/security"

	"gorm.io/gorm"
)

/*
User -> the user itself
*/
type User struct {
	gorm.Model
	Photo    string
	Phone    string
	Email    string          `gorm:"not null;index:usr_email;unique"`
	Password string          `gorm:"not null"`
	Name     string          `gorm:"not null"`
	Surname  string          `gorm:"not null"`
	Birthday time.Time       `gorm:"not null"`
	verified bool            `gorm:"default:false"`
	hasher   security.Hasher `gorm:"-"`
}

func (u *User) checkOrSetHasher() {
	if u.hasher == nil {
		u.hasher = security.NewBcryptHasher()
	}
}

/*
SetPassword -> set user password by hashing the given one. If not hasher is
present, checkOrSetHasher will set the default one.
*/
func (u *User) SetPassword(password string) {
	u.checkOrSetHasher()
	hash, err := u.hasher.GeneratePassword(password)
	if err != nil {
		panic(err)
	}
	u.Password = string(hash)
}

/*
VerifyPassword -> verify if the given password match with the user one. If not
hasher is present, checkOrSetHasher will set the default one.
*/
func (u User) VerifyPassword(password string) bool {
	u.checkOrSetHasher()
	err := u.hasher.VerifyPassword(u.Password, password)
	if err != nil {
		return false
	}
	return true
}
