package controller

import (
	"food-delivery-apps/config"
	"food-delivery-apps/entity"
	"food-delivery-apps/shared"
	"food-delivery-apps/usecase"

	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmployeeController struct{
	menuUc usecase.MenuUseCase
	orderUc usecase.OrderUseCase
	PromoUc usecase.PromoUseCase
	rg *gin.RouterGroup
}

func (c *EmployeeController) Route(){
	c.rg.POST(config.AddMenu, c.AddMenuHandler)
	c.rg.PUT(config.UpdateMenu, c.UpdateMenuHandler)
	c.rg.DELETE(config.DeleteMenu, c.DeleteMenuHandler)
	c.rg.POST(config.AddPromo, c.AddPromoHandler)
	c.rg.GET(config.GetPromo, c.GetPromoHandler)
	c.rg.DELETE(config.DeletePromo, c.DeletePromoHandler)
	c.rg.GET(config.GetAllOrder, c.GetAllOrderHandler)
	c.rg.PATCH(config.UpdateOrderStatus, c.UpdateOrderStatusHandler)
}

// @Summary Create Menu.
// @Description Add a new menu items
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param menuBody body model.MenuRequest true "menu request body"
// @Success 201 {object} model.SingleMenuResponse
// @Failure 400 {object} model.Status "Invalid request payload"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /menu [post]
func (c *EmployeeController) AddMenuHandler(ctx *gin.Context){
	// Retrieve employeeId from JWT auth middleware
	createdBy := ctx.MustGet("userID").(string)

	// Bind JSON request body to menu payload and handle binding errors
	var payload entity.Menu
	if err := ctx.ShouldBindJSON(&payload); err != nil{
		shared.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Set employeeId in payload from JWT data
	payload.CreatedBy = createdBy

	// Call the usecase to create menu
	response, err := c.menuUc.CreateNewMenu(payload)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with created order information
	shared.SendCreateResponse(ctx, response, "successfully created menu")
}

// @Summary Update Menu.
// @Description Update an existing menu.
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Menu ID"
// @Param menuBody body model.MenuRequest true "menu request body"
// @Success 201 {object} model.SingleMenuResponse
// @Failure 400 {object} model.Status "Invalid request payload"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /menu/{id} [put]
func (c *EmployeeController) UpdateMenuHandler(ctx *gin.Context){
	// Extract ID from URL parameter
	id := ctx.Param("id")
	// Retrieve employeeId from JWT auth middleware
	createdBy := ctx.MustGet("userID").(string)
	
	// Bind JSON request body to Menu payload and handle binding errors
	var payload entity.Menu
	if err := ctx.ShouldBindJSON(&payload); err != nil{
		shared.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Set id in payload from URL parameter
	payload.Id = id
	// Set employeeId in payload from JWT data
	payload.CreatedBy = createdBy
	
	// Call the usecae to update specified menu
	resp, err := c.menuUc.UpdateMenu(payload)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with updated order information
	shared.SendSingleResponse(ctx, resp, "successfully updated menu")
}

// @Summary Delete Menu.
// @Description Delete an existing menu.
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Menu ID"
// @Success 204 {object} nil "Successfully deleted menu"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /menu/{id} [delete]
func (c *EmployeeController) DeleteMenuHandler(ctx *gin.Context){
	// Extract ID from URL parameter
	id := ctx.Param("id")
	
	// Call the usecase to delete specified menu
	err := c.menuUc.DeleteMenu(id)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with the provide message
	shared.SendSuccessResponse(ctx, http.StatusNoContent, "successfully deleted menu")
}

// @Summary Create Promo.
// @Description Add a new promo items
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param promoBody body model.PromoRequest true "promo request body"
// @Success 201 {object} model.SinglePromoResponse
// @Failure 400 {object} model.Status "Invalid request payload"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /promo [post]
func (c *EmployeeController) AddPromoHandler(ctx *gin.Context){
	// Retrieve employeeId from JWT auth middleware
	employeeId := ctx.MustGet("userID").(string)
	
	// Bind JSON request body to PromoRequest payload and handle binding errors
	var payload entity.PromoRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil{
		shared.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Set employeeId in payload from JWT data
	payload.EmployeeId = employeeId

	// Call the usecase to create promo
	response, err := c.PromoUc.CreatePromo(payload)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with created promo information
	shared.SendCreateResponse(ctx, response, "successfully created promo")
}

// @Summary Get Promo.
// @Description Retrieves a paginated list of all promo.
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of items per page" default(10)
// @Success 200 {object} model.PagedPromoResponse "Successfully retrieved promo"
// @Failure 404 {object} model.Status "Promo not found"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /promo [get]
func (c *EmployeeController) GetPromoHandler(ctx *gin.Context){
	// Set default pagination parameters (page and size)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	// Call the usecase to fetch promos and pagination info
	resp, paging, err := c.PromoUc.GetAllPromo(page, size)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert promo response data to a slice of empty interfaces for generic handling
	var interfaceSlice = make([]interface{}, len(resp))
	for i, v := range resp{
		interfaceSlice[i] = v
	}

	// Check if the promo data is empty, and if so, send a 404 Not Found response
	if len(interfaceSlice) == 0{
		shared.SendErrorResponse(ctx, http.StatusNotFound, "promo data is empty")
		return
	}

	// Send paged response with promo data and pagination details
	shared.SendPagedResponse(ctx, interfaceSlice, paging, "successfully retrieved all promos")
}

// @Summary Delete Promo.
// @Description Delete an existing promo.
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Promo ID"
// @Success 204 {object} nil "Successfully deleted promo"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /promo/{id} [delete]
func (c *EmployeeController) DeletePromoHandler(ctx *gin.Context){
	// Extract ID from URL parameter
	id := ctx.Param("id")
	
	// Call the usecase to delete specified promo
	err := c.PromoUc.DeletePromo(id)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with the provide message
	shared.SendSuccessResponse(ctx, http.StatusNoContent, "successfully deleted promo")
}

// @Summary Get Order.
// @Description Retrieves a paginated list of all customer's order. filter status with 'finish' or 'unfinish'
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of items per page" default(10)
// @Param status query string false "order status filter"
// @Success 200 {object} model.PagedOrderResponse "Successfully retrieved order"
// @Failure 404 {object} model.Status "Order not found"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /order [get]
func (c *EmployeeController) GetAllOrderHandler(ctx *gin.Context){
	// Set default pagination parameters (page and size)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	status := ctx.DefaultQuery("status", "all")

	// Call the usecase to Fetch all order with pagination
	resp, paging, err := c.orderUc.GetAllOrder(page, size, status)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert order response data to a slice of empty interfaces for generic handling
	var interfaceSlice = make([]interface{}, len(resp))
	for i, v := range resp{
		interfaceSlice[i] = v
	}

	// Check if the order data is empty, and if so, send a 404 Not Found response
	if len(interfaceSlice) == 0 {
		shared.SendErrorResponse(ctx, http.StatusNotFound, fmt.Sprintf("no %s orders found", status))
		return
	}

	// Send paged response with order data and pagination details
	shared.SendPagedResponse(ctx, interfaceSlice, paging, fmt.Sprintf("successfully retrieved %s orders", status))
}

// @Summary Update Order Status.
// @Description Update an existing customer's order status.
// @Tags employee
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Order ID"
// @Success 201 {object} model.SingleOrderResponse
// @Failure 400 {object} model.Status "Invalid request payload"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /order-status/{id} [patch]
func (c *EmployeeController) UpdateOrderStatusHandler(ctx *gin.Context){
	// Extract ID from URL parameter
	id := ctx.Param("id")

	// Bind JSON request body to OrderResponse payload and handle binding errors
	var payload entity.OrderResponse
	if err := ctx.ShouldBind(&payload); err != nil{
		shared.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Set id in payload from URL parameter
	payload.Id = id

	// Call the usecase to update the order status by id
	resp, err := c.orderUc.UpdateOrderStatus(payload)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with updated order information
	shared.SendSingleResponse(ctx, resp, "successfully updated order status")
}

func NewEmployeeController(menuUc usecase.MenuUseCase, orderUc usecase.OrderUseCase, promoUc usecase.PromoUseCase, rg *gin.RouterGroup) *EmployeeController{
	return &EmployeeController{menuUc: menuUc, orderUc: orderUc, PromoUc: promoUc, rg: rg}
}