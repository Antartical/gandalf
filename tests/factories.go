package tests

import (
	"fmt"
	"gandalf/models"

	"syreclabs.com/go/faker"
)

// Creates an user with fake data. It won't be saved into the db
func UserFactory() models.User {
	phone := fmt.Sprintf(
		"+%s%s",
		faker.PhoneNumber().AreaCode(),
		faker.PhoneNumber().SubscriberNumber(9),
	)
	return models.NewUser(
		faker.Internet().Email(),
		faker.Internet().Password(10, 14),
		faker.Name().FirstName(),
		faker.Name().LastName(),
		faker.Date().Birthday(18, 34),
		phone,
	)
}

// Creates an app with fake data. It won't be saved into the db
func AppFactory() models.App {
	user := UserFactory()
	app := models.NewApp(
		faker.Company().Name(),
		faker.Internet().Url(),
		[]string{faker.Internet().Url()},
		user,
	)
	app.User = user
	return app
}
