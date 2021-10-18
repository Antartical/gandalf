package tests

import (
	"gandalf/models"

	"syreclabs.com/go/faker"
)

/*
UserFactory -> user's factory
*/
func UserFactory() models.User {
	return models.NewUser(
		faker.Internet().Email(),
		faker.Internet().Password(8, 14),
		faker.Name().FirstName(),
		faker.Name().LastName(),
		faker.Date().Birthday(18, 34),
		faker.PhoneNumber().CellPhone(),
	)
}

/*
AppFactory -> app's factory
*/
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
