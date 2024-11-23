package controller

import (
	"food-delivery-apps/config"
	"food-delivery-apps/entity"
	"food-delivery-apps/shared"
	"food-delivery-apps/usecase"

	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	uc usecase.UserUseCase
	rg *gin.RouterGroup
}

func (c *UserController) Route() {
	c.rg.PUT(config.UpdateUser, c.UpdateUserHandler)
	c.rg.POST(config.Logout, c.LogoutHandler)
}

// @Summary Update User.
// @Description Update an existing user's own data.
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param userBody body model.UserRequest true "user request body"
// @Success 201 {object} model.SingleUserResponse
// @Failure 400 {object} model.Status "Invalid request payload"
// @Failure 500 {object} model.Status "Internal server error"
// @Security BearerAuth
// @Router /user [put]
func (c *UserController) UpdateUserHandler(ctx *gin.Context){
	// Retrieve userId from JWT auth middleware
	id := ctx.MustGet("userID").(string)

	// Bind JSON request body to User payload and handle binding errors
	var payload entity.User
	if err := ctx.ShouldBindJSON(&payload); err != nil{
		shared.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Set userId in payload from JWT data
	payload.Id = id

	// Call the usecase to update specified user
	updatedUser, err := c.uc.UpdateUser(payload)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with updated order information
	shared.SendSingleResponse(ctx, updatedUser, "successfully updated user")
}

// @Summary Logout User.
// @Description Logout to prevent user access the resource.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} nil "User Logged out successfully"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Security BearerAuth
// @Router /auth/logout [post]
func (c *UserController) LogoutHandler(ctx *gin.Context) {
	// Retrieve token from JWT auth middleware
	tokenString := ctx.MustGet("tokenString").(string)

	// Call the use case to log out
	err := c.uc.Logout(tokenString)
	if err != nil {
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with the provide message
	shared.SendSuccessResponse(ctx, http.StatusOK, "User logged out successfully")
}



func NewUserController(uc usecase.UserUseCase, rg *gin.RouterGroup) *UserController{
	return &UserController{uc: uc, rg: rg}
}