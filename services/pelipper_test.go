package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"gandalf/validators"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type postRecorder struct {
	url         string
	contentType string
	body        map[string]json.RawMessage
}

type mockPost struct {
	postRecorder *postRecorder

	returnCode  int
	raisedError error
}

func newMockPost(code int, err error) mockPost {
	return mockPost{
		postRecorder: new(postRecorder),
		returnCode:   code,
		raisedError:  err,
	}
}

func (mock *mockPost) post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	var payload []byte
	var mappedBody map[string]json.RawMessage
	body.Read(payload)
	json.Unmarshal(payload, &mappedBody)

	mock.postRecorder.url = url
	mock.postRecorder.contentType = contentType
	mock.postRecorder.body = mappedBody

	mockResponse := &http.Response{
		StatusCode: mock.returnCode,
	}
	return mockResponse, mock.raisedError
}

func TestPelipperService(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		pelipperService := NewPelipperService()

		assert.Equal(pelipperService.Host, os.Getenv("PELIPPER_HOST"))
		assert.Equal(pelipperService.SMPTAccount, os.Getenv("PELIPPER_SMTP_ACCOUNT"))
	})

	t.Run("Test manageResponse", func(t *testing.T) {
		raisedError := errors.New("wrong")
		email := "test@test.com"
		mockPost := newMockPost(http.StatusCreated, nil)
		pelipperService := PelipperService{
			Host:        "",
			SMPTAccount: "",
			post:        mockPost.post,
		}

		pelipperService.manageResponse(nil, raisedError, email)
	})

	t.Run("Test SendUserVerifyEmail successfully", func(t *testing.T) {
		host := "miscohost"
		expectedURL := fmt.Sprintf("%s/emails/users/verify", host)
		smtpAccount := "miscoAccount"
		email := "test@test.com"
		mockPost := newMockPost(http.StatusCreated, nil)
		pelipperService := PelipperService{
			Host:        host,
			SMPTAccount: smtpAccount,
			post:        mockPost.post,
		}

		emailData := validators.PelipperUserVerifyEmail{
			Email:            email,
			Name:             "",
			Subject:          "",
			VerificationLink: "",
		}

		pelipperService.SendUserVerifyEmail(emailData)
		assert.Equal(mockPost.postRecorder.url, expectedURL)
		assert.Equal(mockPost.postRecorder.contentType, "application/json")
	})

	t.Run("Test SendUserChangePasswordEmail successfully", func(t *testing.T) {
		host := "miscohost"
		expectedURL := fmt.Sprintf("%s/emails/users/change_password", host)
		smtpAccount := "miscoAccount"
		email := "test@test.com"
		mockPost := newMockPost(http.StatusCreated, nil)
		pelipperService := PelipperService{
			Host:        host,
			SMPTAccount: smtpAccount,
			post:        mockPost.post,
		}

		emailData := validators.PelipperUserChangePassword{
			Email:              email,
			Name:               "",
			Subject:            "",
			ChangePasswordLink: "",
		}

		pelipperService.SendUserChangePasswordEmail(emailData)
		assert.Equal(mockPost.postRecorder.url, expectedURL)
		assert.Equal(mockPost.postRecorder.contentType, "application/json")
	})

}
