package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gandalf/validators"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

// Pelipper service interface
type IPelipperService interface {
	SendUserVerifyEmail(data validators.PelipperUserVerifyEmail)
	SendUserChangePasswordEmail(data validators.PelipperUserChangePassword)
}

// Pelipper is a service through the one we can send notifications to users
type PelipperService struct {
	Host        string
	SMPTAccount string

	post func(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

// Creates a new pelipper service
func NewPelipperService() PelipperService {
	return PelipperService{
		Host:        os.Getenv("PELIPPER_HOST"),
		SMPTAccount: os.Getenv("PELIPPER_SMTP_ACCOUNT"),
		post:        http.Post,
	}
}

func (service PelipperService) manageResponse(response *http.Response, err error, email string) {
	if err != nil || response.StatusCode != http.StatusCreated {
		log.Println(fmt.Sprintf("%s -> %s", err.Error(), service.Host))
		log.Println(fmt.Sprintf("Verification email cannot be sended to %s", email))
	}
}

// Sends the verification email
func (service PelipperService) SendUserVerifyEmail(data validators.PelipperUserVerifyEmail) {
	payload, _ := json.Marshal(map[string]string{
		"from":              service.SMPTAccount,
		"to":                data.Email,
		"name":              data.Name,
		"subject":           data.Subject,
		"verification_link": data.VerificationLink,
	})
	httptest.NewRecorder()

	response, err := service.post(fmt.Sprintf("%s/emails/users/verify", service.Host), "application/json", bytes.NewBuffer(payload))
	service.manageResponse(response, err, data.Email)
}

// Sends the verification email
func (service PelipperService) SendUserChangePasswordEmail(data validators.PelipperUserChangePassword) {
	payload, _ := json.Marshal(map[string]string{
		"from":                 service.SMPTAccount,
		"to":                   data.Email,
		"name":                 data.Name,
		"subject":              data.Subject,
		"change_password_link": data.ChangePasswordLink,
	})
	httptest.NewRecorder()

	response, err := service.post(fmt.Sprintf("%s/emails/users/change_password", service.Host), "application/json", bytes.NewBuffer(payload))
	service.manageResponse(response, err, data.Email)
}
