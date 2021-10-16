package main

import (
	"fmt"
	"gandalf/connections"
	"gandalf/models"
	"gandalf/services"
	"gandalf/validators"
	"os"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mitchellh/mapstructure"
	"github.com/thatisuday/commando"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		AddFlag("email,e", "user email", commando.String, nil).                                                                                     // required
		AddFlag("password,p", "user password", commando.String, nil).                                                                               // required
		AddFlag("name,n", "user name", commando.String, "Agapito").                                                                                 // required
		AddFlag("surname,sn", "user surname", commando.String, "Disousa").                                                                          // required
		AddFlag("birthday,b", "user birthday (must have the following format: YYYY-MM-DDT00:00:00Z)", commando.String, "2006-01-02T15:04:05.000Z"). // required
		AddFlag("phone,b", "user email", commando.String, "+34666123456").                                                                          // required
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			layout := "2006-01-02T15:04:05.000Z"
			birthday, err := time.Parse(layout, fmt.Sprintf("%v", flags["birthday"].Value))
			if err != nil {
				fmt.Println("Birthday has incorrect format")
				os.Exit(1)
			}
			data := map[string]interface{}{
				"Email":           flags["email"].Value,
				"Password":        flags["password"].Value,
				"Name":            flags["name"].Value,
				"Surname":         flags["surname"].Value,
				"Birthday":        birthday,
				"Phone":           flags["phone"].Value,
				"VerificationURL": "",
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

	commando.
		Register("msql").
		SetShortDescription("Inspect the SQL generated by a new migration.").
		SetDescription("Inspect the SQL generated by a new migration.").
		AddFlag("model,m", "model", commando.String, nil). // required                                                                       // required
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			modelName := flags["model"].Value
			_models := map[interface{}]interface{}{
				"App":  &models.App{},
				"User": &models.User{},
			}
			model := _models[modelName]

			if model == nil {
				fmt.Printf("Model: %s does not exists\n", modelName)
				os.Exit(1)
			}

			sqlDB, _, err := sqlmock.New()
			if err != nil {
				panic(err)
			}

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: sqlDB,
			}), &gorm.Config{})
			if err != nil {
				panic(err)
			}

			defer sqlDB.Close()
			gormDB.Migrator().CreateTable(model)

			// func (r *RecorderLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error){
			// 	sql, _ := fc()
			// 	r.Statements = append(r.Statements, sql)
			// }

			// recorder := RecorderLogger{logger: logger.Default.LogMode(logger.Info)}
			// db := connections.NewGormPostgresConnection().Connect()
			// session := db.Session(&gorm.Session{
			// 	Logger: &recorder,
			// })
			// session.AutoMigrate(model)

			// newMigration := sqlize.NewSqlize(sqlize.WithSqlTag("sql"), sqlize.WithMigrationFolder("./migrations"), sqlize.WithPostgresql())
			// _ = newMigration.FromObjects(models.User{})
			// fmt.Printf("STATEMENT: %s", newMigration.StringUp())

			// db := connections.NewGormPostgresConnection().Connect()
			// db.AutoMigrate(&models.User{})
			// db.AutoMigrate(&models.App{})

		})

	// parse command-line arguments
	commando.Parse(nil)

}
