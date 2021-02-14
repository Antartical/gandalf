package security

import "golang.org/x/crypto/bcrypt"

/*
Hasher -> interface for user password hasher
*/
type Hasher interface {
	GeneratePassword(password string) ([]byte, error)
	VerifyPassword(hashedPassword string, plainPassword string) error
}

/*
BcryptHasher -> password hasher based on bcrypt module
*/
type BcryptHasher struct {
	generateFromPassword   func(password []byte, cost int) ([]byte, error)
	compareHashAndPassword func(hashedPassword []byte, password []byte) error
}

/*
GeneratePassword -> generate password by hashing it with bcrypt module
*/
func (hasher BcryptHasher) GeneratePassword(password string) ([]byte, error) {
	return hasher.generateFromPassword([]byte(password), bcrypt.DefaultCost)
}

/*
VerifyPassword -> verifies if the given plainpassword match with the
hashed one with bcrypt module
*/
func (hasher BcryptHasher) VerifyPassword(hashedPassword string, plainPassword string) error {
	return hasher.compareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

/*
NewBcryptHasher -> creates a new BcryptHasher
*/
func NewBcryptHasher() BcryptHasher {
	return BcryptHasher{
		generateFromPassword:   bcrypt.GenerateFromPassword,
		compareHashAndPassword: bcrypt.CompareHashAndPassword,
	}
}
