package validators

import "time"

/*
UserCreateData -> user data for creation
*/
type UserCreateData struct {
	Email    string    `json:"email" binding:"required,email"`
	Password string    `json:"password" binding:"required,alphanumunicode,min=10"`
	Name     string    `json:"name" binding:"required,alphaunicode"`
	Surname  string    `json:"surname" binding:"required,alphaunicode"`
	Birthday time.Time `json:"birthday" binding:"required"`
	Phone    string    `json:"phone" binding:"omitempty,e164"`
}

/*
UserUpdateData -> user data for update
*/
type UserUpdateData struct {
	Password string `json:"password" binding:"omitempty,alphanumunicode,min=10"`
	Phone    string `json:"phone" binding:"omitempty,e164"`
}
