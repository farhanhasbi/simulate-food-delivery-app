package controller

import (
	"food-delivery-apps/config"
	"food-delivery-apps/entity"
	"food-delivery-apps/shared"
	"food-delivery-apps/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerController struct{
	orderUc usecase.OrderUseCase
	balanceUc usecase.BalanceUseCase
	reviewUc usecase.ReviewUseCase
	promoUc usecase.PromoUseCase
	rg *gin.RouterGroup
}

func (c *CustomerController) Route(){
	c.rg.POST(config.CreateBalance, c.CreateBalanceHandler)
	c.rg.GET(config.GetBalance, c.GetBalanceDataHandler)
	c.rg.GET(config.GetPromoCust, c.GetPromoForCustomerHandler)
	c.rg.POST(config.AddOrder, c.AddOrderHandler)
	c.rg.GET(config.GetUnfinishOrder, c.GetUnfinishOrderHandler)
	c.rg.GET(config.GetOrderHistory, c.GetOrderHistoryHandler)
	c.rg.POST(config.AddReview, c.AddReviewHandler)
	c.rg.PUT(config.UpdateReview, c.UpdateReviewHandler)
	c.rg.DELETE(config.DeleteReview, c.DeleteReviewHandler)
}

// @Summary Create Customer's Balance.
// @Description Add amount to increase balance for specific customer
// @Tags customer
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param balanceBody body model.BalanceRequest true "balance request body"
// @Success 201 {object} model.SingleBalanceResponse
// @Failure 400 {object} model.Status "Invalid request payload"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /balance [post]
func (c *CustomerController) CreateBalanceHandler(ctx *gin.Context){
	// Retrieve customerId from JWT auth middleware
	customerId := ctx.MustGet("userID").(string)
	
	// Bind JSON request body to Balance payload and handle binding errors
	var payload entity.Balance
	if err := ctx.ShouldBindJSON(&payload); err != nil{
		shared.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	} 

	// Set customerId in payload from JWT data
	payload.CustomerId = customerId

	// Call the usecase to create balance for specific customer
	resp, err := c.balanceUc.IncreaseBalance(payload)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with created balance information
	shared.SendCreateResponse(ctx, resp, "successfully created balance")
}

// @Summary Get Customer's Balance.
// @Description Retrieves a paginated list of customer's balance.
// @Tags customer
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of items per page" default(10)
// @Success 200 {object} model.PagedBalanceResponse "Successfully retrieved balances"
// @Failure 404 {object} model.Status "Balances not found"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /balance [get]
func (c *CustomerController) GetBalanceDataHandler(ctx*gin.Context){
	// Set default pagination parameters (page and size)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	// Retrieve customerId from JWT auth middleware
	customerId := ctx.MustGet("userID").(string)
	
	// Call the usecase to fetch balances and pagination info for specific customer
	resp, paging, err := c.balanceUc.GetBalanceData(page, size, customerId)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert balance response data to a slice of empty interfaces for generic handling
	var interfaceSlice = make([]interface{}, len(resp))
	for i, v := range resp{
		interfaceSlice[i] = v
	}

	// Check if the balance data is empty, and if so, send a 404 Not Found response
	if len(interfaceSlice) == 0{
		shared.SendErrorResponse(ctx, http.StatusNotFound, "balances not found")
		return
	}

	// Send paged response with balance data and pagination details
	shared.SendPagedResponse(ctx, interfaceSlice, paging, "successfully retrieved balances")
}

// @Summary Get Customer's Promo.
// @Description Retrieves a paginated list of available promo for customer.
// @Tags customer
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
// @Router /available-promo [get]
func (c *CustomerController) GetPromoForCustomerHandler(ctx*gin.Context){
	// Set default pagination parameters (page and size)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	// Retrieve customerId from JWT auth middleware
	customerId := ctx.MustGet("userID").(string)
	
	// Call the usecase to fetch balances and pagination info for specific customer
	resp, paging, err := c.promoUc.GetPromoForCustomer(page, size, customerId)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert balance response data to a slice of empty interfaces for generic handling
	var interfaceSlice = make([]interface{}, len(resp))
	for i, v := range resp{
		interfaceSlice[i] = v
	}

	// Check if the balance data is empty, and if so, send a 404 Not Found response
	if len(interfaceSlice) == 0{
		shared.SendErrorResponse(ctx, http.StatusNotFound, "promo data is empty")
		return
	}

	// Send paged response with balance data and pagination details
	shared.SendPagedResponse(ctx, interfaceSlice, paging, "successfully retrieved promo")
}

// @Summary Create Customer's Order.
// @Description Place new order for specific customer
// @Tags customer
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param orderBody body model.OrderRequest true "order request body"
// @Success 201 {object} model.SingleOrderResponse
// @Failure 400 {object} model.Status "Invalid request payload"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /order [post]
func (c *CustomerController) AddOrderHandler(ctx *gin.Context){
	// Retrieve customerId from JWT auth middleware
	customerId := ctx.MustGet("userID").(string)
	
	// Bind JSON request body to Order payload and handle binding errors
	var payload entity.Order
	if err := ctx.ShouldBind(&payload); err != nil{
		shared.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Set customerId in payload from JWT data
	payload.CustomerId = customerId

	// Call the usecase to create order for specific customer
	resp, err := c.orderUc.CreateNewOrder(payload)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	
	// Send successfully response with created order information
	shared.SendCreateResponse(ctx, resp, "successfully created order")
}

// @Summary Get Unfinish Customer's Order.
// @Description Retrieves unfinish customer's order to track the order status.
// @Tags customer
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} model.SingleOrderResponse "Successfully retrieved customer's order"
// @Failure 404 {object} model.Status "unfinish order not found"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /unfinish-order [get]
func (c *CustomerController) GetUnfinishOrderHandler(ctx*gin.Context){
	// Retrieve customerId from JWT auth middleware
	customerId := ctx.MustGet("userID").(string)
	
	// Call the usecase to retrieve the unfinish order for specific customer
	resp, err := c.orderUc.GetUnfinishCustomerOrder(customerId)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send Succesfully response with unfinish customer's order data
	shared.SendSingleResponse(ctx, resp, "successfully retrieved customer's order")
}

// @Summary Get Finish Customer's Order.
// @Description Retrieves a paginated list of customer's order that already delivered.
// @Tags customer
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Number of items per page" default(10)
// @Param startDate query string false "Start date filter in YYYY-MM-DD format"
// @Param endDate query string false "End date filter in YYYY-MM-DD format"
// @Success 200 {object} model.PagedOrderResponse "Successfully retrieved customer's order"
// @Failure 404 {object} model.Status "Order history not found"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /finish-order [get]
func (c *CustomerController) GetOrderHistoryHandler(ctx *gin.Context){
	// Set default pagination parameters (page and size)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	// Retrieve optional startDate end endDate filter from query
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	// Set customerId in payload from JWT data
	customerId := ctx.MustGet("userID").(string)
	
	// Call the usecase to retrieve the finish order for specific customer
	resp, paging, err := c.orderUc.GetOrderHistory(page, size, startDate, endDate, customerId)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert finish customer's order response data to a slice of empty interfaces for generic handling
	var interfaceSlice = make([]interface{}, len(resp))
	for i, v := range resp{
		interfaceSlice[i] = v
	}

	// Check if the finish customer's order data is empty, and if so, send a 404 Not Found response	
	if len(interfaceSlice) == 0{
		shared.SendErrorResponse(ctx, http.StatusNotFound, "order history not found")
		return
	}

	// Send paged response with finish customer's order data and pagination details
	shared.SendPagedResponse(ctx, interfaceSlice, paging, "successfully retrieved customer's order")
}

// @Summary Create Review.
// @Description add review (1-5) for specific menu by specific customer.
// @Tags customer
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param reviewBody body model.CreateReviewRequest true "review request body"
// @Success 201 {object} model.SingleReviewResponse
// @Failure 400 {object} model.Status "Invalid request payload"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /review [post]
func (c *CustomerController) AddReviewHandler(ctx *gin.Context){
	// Retrieve customerId from JWT auth middleware
	customerId := ctx.MustGet("userID").(string)

	// Bind JSON request body to Review payload and handle binding errors
	var payload entity.Review
	if err := ctx.ShouldBindJSON(&payload); err != nil{
		shared.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Set customerId in payload from JWT data
	payload.CustomerId = customerId

	// Call the usecase to create review by specific customer
	response, err := c.reviewUc.AddReview(payload)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with created review information
	shared.SendCreateResponse(ctx, response, "successfully created review")
}

// @Summary Update Review.
// @Description Update an existing customer's review.
// @Tags customer
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Review ID"
// @Param reviewBody body model.UpdateReviewRequest true "review request body"
// @Success 201 {object} model.SingleReviewResponse
// @Failure 400 {object} model.Status "Invalid request payload"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /review/{id} [put]
func (c *CustomerController) UpdateReviewHandler(ctx*gin.Context){
	// Extract ID from URL parameter
	id := ctx.Param("id")
	// Retrieve customerId from JWT auth middleware
	customerId := ctx.MustGet("userID").(string)
	
	// Bind JSON request body to Review payload and handle binding errors
	var payload entity.Review
	if err := ctx.ShouldBindJSON(&payload); err != nil{
		shared.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	
	// Set id in payload from URL parameter
	payload.Id = id
	// Set customerId in payload from JWT data
	payload.CustomerId = customerId

	// Call the usecase to update review for specific customer
	resp, err := c.reviewUc.UpdateReview(payload)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send successfully response with updated review information
	shared.SendSingleResponse(ctx, resp, "successfully updated review")
}

// @Summary Delete Review.
// @Description Delete an existing customer's review.
// @Tags customer
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Review ID"
// @Success 204 {object} nil "Successfully deleted review"
// @Failure 500 {object} model.Status "Internal server error"
// @Failure 401 {object} model.Status "Unauthorized"
// @Failure 403 {object} model.Status "Forbidden access"
// @Security BearerAuth
// @Router /review/{id} [delete]
func (c *CustomerController) DeleteReviewHandler(ctx*gin.Context){
	// Extract ID from URL parameter
	id := ctx.Param("id")
	// Set customerId in payload from JWT data
	customerId := ctx.MustGet("userID").(string)
	
	// Call the usecase to delete specified review
	err := c.reviewUc.DeleteReview(id, customerId)
	if err != nil{
		shared.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Send a successfully response with the provided message
	shared.SendSuccessResponse(ctx, http.StatusNoContent, "successfully deleted review")
}

func NewCustomerController(orderUc usecase.OrderUseCase, balanceUc usecase.BalanceUseCase, reviewUc usecase.ReviewUseCase, promoUc usecase.PromoUseCase, rg *gin.RouterGroup) *CustomerController{
	return &CustomerController{orderUc: orderUc, balanceUc: balanceUc, reviewUc: reviewUc, promoUc: promoUc, rg: rg}
}