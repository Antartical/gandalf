package services

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type userServiceWithMockDB struct {
	service UserService
	sqlmock sqlmock.Sqlmock
}

// func newUserServiceWIthMockDB() userServiceWithMockDB {
// 	sqldb, mock, err := sqlmock.New()
// 	panic(err)

// 	db, err := gorm.Open(postgres.New(postgres.Config{
// 		Conn: sqldb,
// 	}), &gorm.Config{})
// 	panic(err)

// 	return userServiceWithMockDB{
// 		service: UserService{db},
// 		sqlmock: mock,
// 	}
// }

func TestUserServiceCreate(t *testing.T) {
	//assert := require.New(t)

	t.Run("Test user create successfully", func(t *testing.T) {

	})

	t.Run("Test user create database error", func(t *testing.T) {

	})
}
