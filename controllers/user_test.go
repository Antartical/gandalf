package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"gandalf/models"
	"gandalf/services"
	"gandalf/validators"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

type createRecorder struct {
	userData validators.UserCreateData
}

type uuidRecorder struct {
	uuid uuid.UUID
}

type updateRecorder struct {
	uuid     uuid.UUID
	userData validators.UserUpdateData
}

type mockUserService struct {
	createRecorder *createRecorder
	readRecorder   *uuidRecorder
	updateRecorder *updateRecorder
	deleteRecorder *uuidRecorder

	createError error
	readError   error
	updateError error
	deleteError error
}

func (service *mockUserService) Create(userData validators.UserCreateData) (models.User, error) {
	*service.createRecorder = createRecorder{userData: userData}
	return models.User{}, service.createError
}

func (service *mockUserService) Read(uuid uuid.UUID) (models.User, error) {
	*service.readRecorder = uuidRecorder{uuid: uuid}
	return models.User{}, service.readError
}

func (service *mockUserService) Update(uuid uuid.UUID, userData validators.UserUpdateData) (models.User, error) {
	*service.updateRecorder = updateRecorder{uuid: uuid, userData: userData}
	return models.User{}, service.updateError
}

func (service *mockUserService) Delete(uuid uuid.UUID) error {
	*service.deleteRecorder = uuidRecorder{uuid: uuid}
	return service.deleteError
}

func newMockedService(createError error, readError error, updateError error, deleteError error) mockUserService {
	return mockUserService{
		createRecorder: new(createRecorder),
		readRecorder:   new(uuidRecorder),
		updateRecorder: new(updateRecorder),
		deleteRecorder: new(uuidRecorder),
		createError:    createError,
		readError:      readError,
		updateError:    updateError,
		deleteError:    deleteError,
	}
}

func setupUserRouter(userService services.IUserService) *gin.Engine {
	router := gin.Default()
	RegisterUserRoutes(router, userService)
	return router
}

func TestCreateUser(t *testing.T) {
	assert := require.New(t)

	t.Run("Test create user successfully", func(t *testing.T) {
		userService := newMockedService(nil, nil, nil, nil)
		router := setupUserRouter(&userService)
		var response gin.H

		email := "test@test.com"
		password := "testtesttesttest"
		name := "test"
		surname := "test"
		birthdayStr := "2021-02-15T00:00:00Z"
		birthdayDat, _ := time.Parse(time.RFC3339, birthdayStr)
		payload, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
			"name":     name,
			"Surname":  surname,
			"Birthday": birthdayStr,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusCreated)
		assert.Equal(userService.createRecorder.userData.Email, email)
		assert.Equal(userService.createRecorder.userData.Password, password)
		assert.Equal(userService.createRecorder.userData.Name, name)
		assert.Equal(userService.createRecorder.userData.Surname, surname)
		assert.Equal(userService.createRecorder.userData.Birthday, birthdayDat)
	})

	t.Run("Test create user binding error", func(t *testing.T) {
		userService := newMockedService(nil, nil, nil, nil)
		router := setupUserRouter(&userService)
		var response gin.H

		payload, _ := json.Marshal(map[string]string{
			"email": "wrong email",
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("Test create user service error", func(t *testing.T) {
		expectedError := errors.New("create error")
		userService := newMockedService(expectedError, nil, nil, nil)
		router := setupUserRouter(&userService)
		var response gin.H

		email := "test@test.com"
		password := "testtesttesttest"
		name := "test"
		surname := "test"
		birthdayStr := "2021-02-15T00:00:00Z"
		payload, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
			"name":     name,
			"Surname":  surname,
			"Birthday": birthdayStr,
		})

		recorder := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
		router.ServeHTTP(recorder, request)
		json.Unmarshal(recorder.Body.Bytes(), &response)

		assert.Equal(recorder.Result().StatusCode, http.StatusBadRequest)
		assert.Equal(response["error"], expectedError.Error())
	})
}
