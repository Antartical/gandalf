package main

import (
	"fmt"
	"gandalf/bindings"
	"gandalf/connections"
	"gandalf/services"
	"gandalf/validators"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/thatisuday/commando"
	"gorm.io/gorm/logger"
)

type RecorderLogger struct {
	logger.Interface
	Statements []string
}

func main() {

	// configure commando
	commando.
		SetExecutableName("gandalf-cli").
		SetVersion("1.0.0").
		SetDescription("A gandalf cli")

	// configure info command
	commando.
		Register("create-user").
		SetShortDescription("Creates a new user into gandalf database").
		SetDescription("Insert the given user into the gandalf database").
		AddFlag("email,e", "user email", commando.String, nil).                                                             // required
		AddFlag("password,p", "user password", commando.String, nil).                                                       // required
		AddFlag("name,n", "user name", commando.String, "Agapito").                                                         // required
		AddFlag("surname,sn", "user surname", commando.String, "Disousa").                                                  // required
		AddFlag("birthday,b", "user birthday (must have the following format: YYYY-MM-DD)", commando.String, "1997-12-21"). // required
		AddFlag("phone,ph", "user email", commando.String, "+34666123456").                                                 // required
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			var birthdate bindings.BirthDate
			if err := birthdate.UnmarshalJSON([]byte(fmt.Sprintf("%v", flags["birthday"].Value))); err != nil {
				fmt.Println("Birthday has incorrect format")
				os.Exit(1)
			}
			data := map[string]interface{}{
				"Email":    flags["email"].Value,
				"Password": flags["password"].Value,
				"Name":     flags["name"].Value,
				"Surname":  flags["surname"].Value,
				"Birthday": birthdate,
				"Phone":    flags["phone"].Value,
			}

			var input validators.UserCreateData
			if err := mapstructure.Decode(data, &input); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			db := connections.NewGormPostgresConnection().Connect()
			userService := services.NewUserService(db)
			user, err := userService.Create(input)
			if err != nil {
				fmt.Printf("User %s already exists\n", input.Name)
				os.Exit(1)
			}

			userService.Verificate(user)
			fmt.Printf("User %s created successfully\n", user.Name)
		})

	// parse command-line arguments
	commando.Parse(nil)

}
