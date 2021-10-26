package controllers

import (
	"fmt"
	"gandalf/helpers"
	"gandalf/middlewares"
	"gandalf/security"
	"gandalf/serializers"
	"gandalf/services"
	"gandalf/validators"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// Register user endpoints to the given router
func RegisterUserRoutes(
	router *gin.Engine,
	authBearerMiddleware middlewares.IAuthBearerMiddleware,
	authService services.IAuthService,
	userService services.IUserService,
	pelipperService services.IPelipperService,
) {
	controller := UserController{
		authService:     authService,
		userService:     userService,
		pelipperService: pelipperService,
		authMiddleware:  authBearerMiddleware,
	}

	publicRoutes := router.Group("/users")
	{
		publicRoutes.POST("", controller.CreateUser)
	}

	readRoutes := router.Group("/users")
	{
		scopes := []string{security.ScopeUserRead}
		readRoutes.Use(authBearerMiddleware.HasScopes(scopes))

		readRoutes.GET(":uuid", controller.ReadUser)
	}
}

// Controller for /users endpoints
type UserController struct {
	authService     services.IAuthService
	userService     services.IUserService
	pelipperService services.IPelipperService
	authMiddleware  middlewares.IAuthBearerMiddleware
}

// @Summary Create User
// @Description Creates a new user
// @ID user-create
// @Tags User
// @Accept json
// @Produce json
// @Param user body validators.UserCreateData true "Creates a new user"
// @Success 201 {object} serializers.UserSerializer
// @Failure 400 {object} helpers.HTTPError
// @Router /users [post]
func (controller UserController) CreateUser(c *gin.Context) {
	var input validators.UserCreateData
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	user, err := controller.userService.Create(input)
	if err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}

	verifyToken := controller.authService.GenerateTokens(
		*user, []string{security.ScopeUserVerify},
	).AccessToken

	url := os.Getenv("EMAIL_VERIFICATION_URL")
	emailData := validators.PelipperUserVerifyEmail{
		Email:   user.Email,
		Name:    user.Name,
		Subject: "Welcome",
		VerificationLink: fmt.Sprintf(
			"%s?code=%s", url, verifyToken,
		),
	}

	go controller.pelipperService.SendUserVerifyEmail(emailData)
	c.JSON(http.StatusCreated, serializers.NewUserSerializer(*user))
}

// @Summary Get user
// @Description get an user by his uuid
// @ID user-read-uuid
// @Tags User
// @Accept json
// @Produce json
// @Param uuid path string true "User uuid"
// @Success 200 {object} serializers.UserSerializer
// @Failure 400 {object} helpers.HTTPError
// @Failure 403 {object} helpers.HTTPError
// @Router /users/{uuid} [get]
func (controller UserController) ReadUser(c *gin.Context) {
	var input validators.UserReadData
	if err := c.ShouldBindUri(&input); err != nil {
		helpers.AbortWithStatus(c, http.StatusBadRequest, err)
		return
	}
	uuid, _ := uuid.FromString(input.UUID)
	user, err := controller.userService.Read(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, serializers.NewUserSerializer(*user))
}
