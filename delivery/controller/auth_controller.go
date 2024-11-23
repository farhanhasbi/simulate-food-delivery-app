package controller

import (
	"food-delivery-apps/config"
	"food-delivery-apps/entity/dto"
	"food-delivery-apps/shared"
	"food-delivery-apps/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	uc usecase.AuthUseCase
	rg *gin.RouterGroup
}


func (c *AuthController) Route() {
	c.rg.POST(config.Register, c.RegisterHandler)
	c.rg.POST(config.Login, c.LoginHandler)
}

// @Summary Register New User
// @Description Registers a new user with username, email, password, and gender.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param RegisterBody body dto.AuthRequestRegister true "Register User"
// @Success 201 {object} model.SingleUserResponse
// @Failure 400 {object} model.Status "Invalid request payload"
// @Failure 500 {object} model.Status "Internal server error"
// @Router /auth/register [post]
func (c *AuthController) RegisterHandler(ctx *gin.Context){
	// Bind JSON request body to AuthRequestRegister payload and handle binding errors
	var payload dto.AuthRequestRegister
	if err := ctx.ShouldBindJSON(&payload); err != nil{
		shared.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Call the usecase to register user
	user, err := c.uc.Register(payload)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with created user information
	shared.SendCreateResponse(ctx, user, "User registered successfully")
}

// @Summary Login User
// @Description Logs in a user with email and password. Returns a JWT token on success.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param LoginBody body dto.AuthRequestLogin true "Login User"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} model.Status "Invalid request payload"
// @Failure 500 {object} model.Status "Internal server error"
// @Router /auth/login [post]
func (c *AuthController) LoginHandler(ctx *gin.Context) {
	// Bind JSON request body to AuthRequestLogin payload and handle binding errors
	var payload dto.AuthRequestLogin
	if err := ctx.ShouldBindJSON(&payload); err != nil{
		shared.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Call the usecase to login and get the token
	token, err := c.uc.Login(payload)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with the provided token and message
	shared.SendSingleResponse(ctx, token, "User logged in successfully")
}

func NewAuthController(uc usecase.AuthUseCase, rg *gin.RouterGroup) *AuthController{
	return &AuthController{uc: uc, rg: rg}
}