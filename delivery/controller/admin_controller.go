package controller

import (
	"food-delivery-apps/config"
	"food-delivery-apps/shared"
	"food-delivery-apps/usecase"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	uc usecase.UserUseCase
	rg *gin.RouterGroup
}

func (c *AdminController) Route() {
	c.rg.GET(config.GetAllUser, c.GetAllUserHandler)
	c.rg.PATCH(config.Role, c.AssignToEmployeeHandler)
	c.rg.DELETE(config.DeleteUser, c.DeleteUserHandler)
}

// @Summary Get Users
// @Description Retrieves a paginated list of users. You can filter by role.
// @Tags Admin
// @Accept json
// @Producejson
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of items per page" default(10)
// @Param role query string false "User role filter"
// @Success 200 {object} model.PagedUserResponse "Successfully retrieved users"
// @Failure 404 {object} model.Status "No users found"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /user [get]
func (c *AdminController) GetAllUserHandler(ctx *gin.Context){
	// Set default pagination parameters (page and size)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	// Retrieve optional role filter from query
	role := ctx.Query("role")

	// call the usecase to fetch users and pagination info
	resp, paging, err := c.uc.GetAllUser(page, size, role)
	if err != nil {
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert user response data to a slice of empty interfaces for generic handling
	var interfaceSlice = make([]interface{}, len(resp))
	for i, v := range resp{
		interfaceSlice[i] = v
	}

	// Check if the user data is empty, and if so, send a 404 Not Found response
	if len(interfaceSlice) == 0{
		shared.SendErrorResponse(ctx, http.StatusNotFound, "user data is empty")
		return
	}
	
	// Send paged response with user data and pagination details
	shared.SendPagedResponse(ctx, interfaceSlice, paging, "successfully retrieved users")
}

// @Summary Update User Role.
// @Description Update user's role from customer to employee.
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "User ID"
// @Success 201 {object} model.SingleUserResponse
// @Failure 400 {object} model.Status "Invalid request payload"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /user/{id}/role [patch]
func (c *AdminController) AssignToEmployeeHandler(ctx *gin.Context) {
	// Extract target user ID from URL parameter
	targetUserID := ctx.Param("id")

	// Ensure the targetUserID is provided in the request URL
	if targetUserID == "" {
		shared.SendErrorResponse(ctx, http.StatusBadRequest, "user ID is required")
		return
	}
	
	// Call the use case to assign the employee role to the specified user
	updatedUser, err := c.uc.AssignToEmployee(targetUserID)
	if err != nil {
		shared.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
    return
  }

	// Send successful response with updated user information
	shared.SendSingleResponse(ctx, updatedUser, "sucessfully updated role")
}

// @Summary Delete User.
// @Description Delete an existing user.
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "User ID"
// @Success 204 {object} nil "Successfully deleted user"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /user/{id} [delete]
func (c *AdminController) DeleteUserHandler(ctx *gin.Context){
	// Extract target user ID from URL parameter
	targetUserID := ctx.Param("id")
	
	// Call the usecase to delete the specified user
	err := c.uc.DeleteUser(targetUserID)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send a successfully response with the provided message
	shared.SendSuccessResponse(ctx, http.StatusNoContent, "successfully deleted user")
}

func NewAdminController(uc usecase.UserUseCase, rg *gin.RouterGroup) *AdminController{
	return &AdminController{uc: uc, rg: rg}
}