package tests

import (
	"gandalf/models"
	"time"

	"github.com/bxcodec/faker/v3"
)

func userFactory() models.User {
	fakeTime, _ := faker.GetDateTimer().Time()
	return models.NewUser(
		faker.Email(),
		faker.Password(),
		faker.Name(),
		faker.LastName(),
		fakeTime.(time.Time),
		faker.Phonenumber(),
	)
}
