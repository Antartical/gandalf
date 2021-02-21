package serializers

import (
	"gandalf/models"
	"time"

	"github.com/gofrs/uuid"
)

/*
UserSerializer -> user serializer for api output
*/
type UserSerializer struct {
	UUID     uuid.UUID `json:"uuid"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Birthday time.Time `json:"birthday"`
	Phone    string    `json:"phone"`
}

/*
NewUserSerializer -> creates a new user serializer and fills it with
the given user data.
*/
func NewUserSerializer(user models.User) UserSerializer {
	return UserSerializer{
		UUID:     user.UUID,
		Email:    user.Email,
		Name:     user.Name,
		Surname:  user.Surname,
		Birthday: user.Birthday,
		Phone:    user.Phone,
	}
}
