package serializers

import (
	"gandalf/models"
	"time"

	"github.com/gofrs/uuid"
)

type userDataSerializer struct {
	UUID     uuid.UUID `json:"uuid" example:"4722679b-5a48-4e85-9084-605e8df610f4"`
	Email    string    `json:"email" example:"test@test.com"`
	Name     string    `json:"name" example:"John"`
	Surname  string    `json:"surname" example:"Doe"`
	Birthday time.Time `json:"birthday" example:""`
	Phone    string    `json:"phone" example:"+34666123456"`
}

// User serialization struct
type UserSerializer struct {
	ObjectType string             `json:"type" example:"user"`
	Data       userDataSerializer `json:"data"`
}

// Creates a new user serializer and fills it with
// the given user data.
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
