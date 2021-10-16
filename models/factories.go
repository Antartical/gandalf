package models

import (
	"syreclabs.com/go/faker"
)

/*
UserFactory -> user's factory
*/
func UserFactory() User {
	return NewUser(
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
func AppFactory() App {
	return NewApp(
		faker.Company().Name(),
		faker.Internet().Url(),
		[]string{faker.Internet().Url()},
		UserFactory(),
	)
}
