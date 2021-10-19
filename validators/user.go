package validators

import (
	"time"
)

// Validator for resend email notification to an user
type UserResendEmail struct {
	Email string `json:"email" binding:"required,email" example:"johndoe@example.com"`
}

// Validator for user creation
type UserCreateData struct {
	Email    string    `json:"email" binding:"required,email" example:"johndoe@example.com"`
	Password string    `json:"password" binding:"required,min=10" example:"My@appPassw0rd"`
	Name     string    `json:"name" binding:"required" example:"John"`
	Surname  string    `json:"surname" binding:"required" example:"Doe"`
	Birthday time.Time `json:"birthday" binding:"required" time_format:"2006-01-02" time_utc:"1"`
	Phone    string    `json:"phone" binding:"omitempty,e164" example:"+34666123456"`
}

// Validator for retrieve user by his uuid
type UserReadData struct {
	UUID string `uri:"uuid" binding:"required,uuid4" example:"4722679b-5a48-4e85-9084-605e8df610f4"`
}

// Validator for reset user password
type UserResetPasswordData struct {
	Password string `json:"password" binding:"min=10,required" example:"My@appPassw0rd"`
}

// Validator for user update
type UserUpdateData struct {
	Password string `json:"password" binding:"omitempty,min=10" example:"My@appPassw0rd"`
	Phone    string `json:"phone" binding:"omitempty,e164" example:"+34666123456"`
}
