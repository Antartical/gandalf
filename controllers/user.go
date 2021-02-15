package controllers

import (
	"gandalf/serializers"
	"gandalf/services"
	"gandalf/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/*
CreateUser -> creates a new user
*/
func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userService := c.MustGet("userService").(services.IUserService)

	var input validators.UserCreateData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := userService.Create(db, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, serializers.NewUserSerializer(user))
}
