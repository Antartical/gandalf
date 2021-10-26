package serializers

import (
	"gandalf/helpers"
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

type appPublicDataSerializer struct {
	UUID    uuid.UUID `json:"uuid" example:"4722679b-5a48-4e85-9084-605e8df610f4"`
	Name    string    `json:"name" example:"MyApp"`
	IconUrl string    `json:"icon_url" example:"https://rb.gy/1akgfo"`
}

// App serialization struct
type AppSerializer struct {
	ObjectType string            `json:"type" example:"app"`
	Data       appDataSerializer `json:"data"`
}

type paginatedAppsSerializerMeta struct {
	Cursor CursorSerializer `json:"cursor"`
}

type PaginatedAppsSerializer struct {
	ObjectType string                      `json:"type" example:"app"`
	Data       []appDataSerializer         `json:"data"`
	Meta       paginatedAppsSerializerMeta `json:"meta"`
}

type PaginatedAppsPublicSerializer struct {
	ObjectType string                      `json:"type" example:"app"`
	Data       []appPublicDataSerializer   `json:"data"`
	Meta       paginatedAppsSerializerMeta `json:"meta"`
}

// Creates a new app serializer and fills it with
// the given user data.
func NewAppSerializer(app models.App) AppSerializer {
	return AppSerializer{
		ObjectType: "app",
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
func NewPaginatedAppsSerializer(apps []models.App, cursor helpers.Cursor) PaginatedAppsSerializer {
	var serializedApps []appDataSerializer
	for _, app := range apps {
		serializedApp := appDataSerializer{
			UUID:         app.UUID,
			ClientID:     app.ClientID,
			ClientSecret: app.ClientSecret,
			Name:         app.Name,
			IconUrl:      app.IconUrl,
			RedirectUrls: app.RedirectUrls,
		}
		serializedApps = append(serializedApps, serializedApp)
	}

	return PaginatedAppsSerializer{
		ObjectType: "app",
		Data:       serializedApps,
		Meta: paginatedAppsSerializerMeta{
			Cursor: NewCursorSerializer(cursor),
		},
	}
}

// Creates a new apps public serializer and fills it with
// the given user data.
func NewPaginatedAppsPublicSerializer(apps []models.App, cursor helpers.Cursor) PaginatedAppsPublicSerializer {
	var serializedApps []appPublicDataSerializer
	for _, app := range apps {
		serializedApp := appPublicDataSerializer{
			UUID:    app.UUID,
			Name:    app.Name,
			IconUrl: app.IconUrl,
		}
		serializedApps = append(serializedApps, serializedApp)
	}

	return PaginatedAppsPublicSerializer{
		ObjectType: "app",
		Data:       serializedApps,
		Meta: paginatedAppsSerializerMeta{
			Cursor: NewCursorSerializer(cursor),
		},
	}
}
