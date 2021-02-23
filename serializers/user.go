package serializers

import (
	"gandalf/models"
	"time"

	"github.com/gofrs/uuid"
)

type userDataSerializer struct {
	UUID     uuid.UUID `json:"uuid"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Birthday time.Time `json:"birthday"`
	Phone    string    `json:"phone"`
}

/*
UserSerializer -> user serializer for api output
*/
type UserSerializer struct {
	ObjectType string             `json:"type"`
	Data       userDataSerializer `json:"data"`
}

/*
NewUserSerializer -> creates a new user serializer and fills it with
the given user data.
*/
func NewUserSerializer(user models.User) UserSerializer {
	return UserSerializer{
		ObjectType: "user",
		Data: userDataSerializer{
			UUID:     user.UUID,
			Email:    user.Email,
			Name:     user.Name,
			Surname:  user.Surname,
			Birthday: user.Birthday,
			Phone:    user.Phone,
		},
	}
}
