package main

import (
	"fmt"
	"gandalf/bindings"
	"gandalf/connections"
	"gandalf/models"
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

		// configure info command
	commando.
		Register("create-app").
		SetShortDescription("Creates a new app into gandalf database").
		SetDescription("Insert the given app into the gandalf database").                                                                             // required
		AddFlag("name,n", "app name", commando.String, "MyNewApp").                                                                                   // required
		AddFlag("iconurl,ic", "icon url", commando.String, "https://www.vhv.rs/dpng/d/409-4097341_penguin-png-pic-penguins-png-transparent-png.png"). // required
		AddFlag("redirecturl,r", "redirect url", commando.String, nil).                                                                               // required                                                // required
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			url, err := flags["redirecturl"].GetString()
			data := map[string]interface{}{
				"name":          flags["name"].Value,
				"icon_url":      flags["iconurl"].Value,
				"redirect_urls": fmt.Sprintf("[\"%s\"]", url),
			}
			var input validators.AppCreateData
			if err := mapstructure.Decode(data, &input); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			input.IconUrl, _ = flags["iconurl"].GetString()
			input.RedirectUrls = []string{url}
			db := connections.NewGormPostgresConnection().Connect()
			appService := services.NewAppService(db)
			var user *models.User
			db.First(&user)
			if user == nil {
				fmt.Print("You need to create an user first\n")
				os.Exit(1)
			}

			app, err := appService.Create(input, *user)
			if err != nil {
				fmt.Print(err)
				os.Exit(1)
			}

			fmt.Printf("Client ID: %s\nClient secret: %s\nRedirect Url: %s", app.ClientID, app.ClientSecret, app.RedirectUrls)
		})

	// parse command-line arguments
	commando.Parse(nil)

}
