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

/*
IPelipperService -> pelipper interface
*/
type IPelipperService interface {
	SendUserVerifyEmail(data validators.PelipperUserVerifyEmail)
}

/*
PelipperService -> pelipper is a service through the one we can
send notifications to users
*/
type PelipperService struct {
	Host        string
	SMPTAccount string

	post func(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

/*
NewPelipperService -> creates a new pelipper service
*/
func NewPelipperService() PelipperService {
	return PelipperService{
		Host:        os.Getenv("PELIPPER_HOST"),
		SMPTAccount: os.Getenv("PELIPPER_SMTP_ACCOUNT"),
		post:        http.Post,
	}
}

/*
SendUserVerifyEmail -> send the verification email
*/
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
	if err != nil || response.StatusCode != http.StatusCreated {
		log.Println(fmt.Sprintf("%s -> %s", err.Error(), service.Host))
		log.Println(fmt.Sprintf("Verification email cannot be sended to %s", data.Email))
	}
}
