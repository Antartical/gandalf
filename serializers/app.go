package serializers

import (
	"gandalf/models"

	"github.com/gofrs/uuid"
)

type appDataSerializer struct {
	UUID         uuid.UUID `json:"uuid" example:"4722679b-5a48-4e85-9084-605e8df610f4"`
	ClientID     uuid.UUID `json:"client_id" example:"4722679b-5a48-4e85-9084-605e8df610f4"`
	ClientSecret string    `json:"client_secret" example:"iuhgf3874tiu34gtwerbguv3iu74"`
	Name         string    `json:"name" example:"MyApp"`
	IconUrl      string    `json:"icon_url" example:"https://rb.gy/1akgfo"`
	RedirectUrls []string  `json:"redirect_urls" example:"http://localhost:/callback"`
}

// App serialization struct
type AppSerializer struct {
	ObjectType string            `json:"type" example:"app"`
	Data       appDataSerializer `json:"data"`
}

// Creates a new app serializer and fills it with
// the given user data.
func NewAppSerializer(app models.App) AppSerializer {
	return AppSerializer{
		ObjectType: "user",
		Data: appDataSerializer{
			UUID:         app.UUID,
			ClientID:     app.ClientID,
			ClientSecret: app.ClientSecret,
			Name:         app.Name,
			IconUrl:      app.IconUrl,
			RedirectUrls: app.RedirectUrls,
		},
	}
}

// Creates a new apps serializer and fills it with
// the given user data.
func NewAppsSerializer(apps []models.App) []AppSerializer {
	var serializedApps []AppSerializer
	for _, app := range apps {
		serializedApps = append(serializedApps, NewAppSerializer(app))
	}
	return serializedApps
}
