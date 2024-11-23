package controller

import (
	"food-delivery-apps/config"
	"food-delivery-apps/shared"
	"food-delivery-apps/usecase"

	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

type PublicController struct{
	menuUc usecase.MenuUseCase
	reviewUc usecase.ReviewUseCase
	rg *gin.RouterGroup
}

func (c *PublicController) Route(){
	c.rg.GET(config.GetMenu, c.GetMenuHandler)
	c.rg.GET(config.GetReview, c.GetReviewHandler)
}


// @Summary Get Menus
// @Description Retrieves a paginated list of menus. You can filter by type or name.
// @Tags Public
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of items per page" default(10)
// @Param type query string false "Menu type filter"
// @Param name query string false "Menu name filter"
// @Success 200 {object} model.PagedMenuResponse "Successfully retrieved menus"
// @Failure 404 {object} model.Status "No menus found"
// @Failure 500 {object} model.Status "Internal server error"
// @Router /menu [get]
func (c *PublicController) GetMenuHandler(ctx *gin.Context){
	// Set default pagination parameters (page and size)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	// Retrieve optional name and type filter from query
	mtype := ctx.Query("type")
	mname := ctx.Query("name")

	// Call the usecase to fetch menus and pagination info
	resp, paging, err := c.menuUc.GetAllMenu(page, size, mtype, mname)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert menu response data to a slice of empty interfaces for generic handling
	var interfaceSlice = make([]interface{}, len(resp))
	for i, v := range resp{
		interfaceSlice[i] = v
	}

	// Check if the review data is empty, and if so, send a 404 Not Found response
	if len(interfaceSlice) == 0{
		shared.SendErrorResponse(ctx, http.StatusNotFound, "no menus found")
		return
	}
	
	// Send paged response with review data and pagination details
	shared.SendPagedResponse(ctx, interfaceSlice, paging, "successfully retrieved menus")
}

// @Summary Get Reviews
// @Description Retrieves a paginated list of reviews.
// @Tags Public
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of items per page" default(10)
// @Success 200 {object} model.PagedReviewResponse "Successfully retrieved reviews"
// @Failure 404 {object} model.Status "No reviews found"
// @Failure 500 {object} model.Status "Internal server error"
// @Router /review [get]
func (c *PublicController) GetReviewHandler(ctx *gin.Context){
	// Set default pagination parameters (page and size)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	// Call the usecase to fetch reviews and pagination info
	resp, paging, err := c.reviewUc.GetReview(page, size)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert review response data to a slice of empty interfaces for generic handling
	var interfaceSlice = make([]interface{}, len(resp))
	for i, v := range resp{
		interfaceSlice[i] = v
	}

	// Check if the review data is empty, and if so, send a 404 Not Found response
	if len(interfaceSlice) == 0{
		shared.SendErrorResponse(ctx, http.StatusNotFound, "no reviews found")
		return
	}

	// Send paged response with review data and pagination details
	shared.SendPagedResponse(ctx, interfaceSlice, paging, "successfully retrieved reviews")
}


func NewPublicController(menuUc usecase.MenuUseCase, reviewUc usecase.ReviewUseCase, rg *gin.RouterGroup) *PublicController{
	return &PublicController{menuUc: menuUc, reviewUc: reviewUc, rg: rg}
}