package serializers

import "gandalf/helpers"

// cursor data serializer
type cursorDataSerializer struct {
	PageSize     int  `json:"size" example:"5"`
	ActualPage   int  `json:"actual" example:"4"`
	PreviousPage *int `json:"previous" binding:"exists" example:"5"`
	NextPage     *int `json:"next" binding:"exists" example:"6"`
	TotalPages   int  `json:"total_pages" example:"20"`
	TotalObjects int  `json:"total_objects" example:"100"`
}

// Cursor serializer
type CursorSerializer struct {
	ObjectType string               `json:"type" example:"cursor"`
	Data       cursorDataSerializer `json:"data"`
}

// Serializes the given cursor
func NewCursorSerializer(cursor helpers.Cursor) CursorSerializer {
	return CursorSerializer{
		ObjectType: "cursor",
		Data: cursorDataSerializer{
			PageSize:     cursor.PageSize,
			ActualPage:   cursor.Page,
			PreviousPage: cursor.PreviousPage,
			NextPage:     cursor.NextPage,
			TotalPages:   cursor.TotalPages,
			TotalObjects: cursor.Total,
		},
	}
}
