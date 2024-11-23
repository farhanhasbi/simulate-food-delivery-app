package usecase

import (
	"fmt"
	"food-delivery-apps/entity"
	"food-delivery-apps/repository"
	"food-delivery-apps/shared/model"
	"time"
)

type orderUseCase struct{
	repo repository.OrderRepository
	menuRepo repository.MenuRepository
	balanceRepo repository.BalanceRepository
	promoRepo repository.PromoRepository
}

type OrderUseCase interface{
	CreateNewOrder(payload entity.Order) (entity.OrderResponse, error)
	GetUnfinishCustomerOrder(customerId string) (entity.OrderResponse, error)
	UpdateOrderStatus(payload entity.OrderResponse) (entity.OrderResponse, error)
	GetAllOrder(page, size int, status string) ([]entity.OrderResponse, model.Paging, error) 
	GetOrderHistory(page, size int, startDate, endDate string, customerId string) ([]entity.OrderResponse, model.Paging, error)
}

func (uc *orderUseCase) CreateNewOrder(payload entity.Order) (entity.OrderResponse, error){
	// Get totalPrice from CalculateTotalPrice method
	totalPrice, err := uc.CalculateTotalPrice(payload)
	if err != nil{
		return entity.OrderResponse{}, err
	}

	// Check for any undelivered orders by the customer; prevent new orders if one exists.
	var count int
	err = uc.repo.CountUnfishOrder(payload.CustomerId, &count)
	if err != nil{
		return entity.OrderResponse{}, err
	}
	if count > 0{
		return entity.OrderResponse{}, fmt.Errorf("cannot place a new order until the current one is delivered")
	}

	// Apply discount if a promo code is provided.
	if payload.PromoCode != ""{
		promo, err := uc.ApplyPromo(payload)
		if err != nil{
			return entity.OrderResponse{}, err
		}

		// Adjust total price based on promo type (percentage or flat).
		if promo.IsPercentage {
			discount := (totalPrice * promo.Discount) / 100
			totalPrice -= discount
		} else {
			totalPrice -= promo.Discount
		}
	}

	// Ensure total price is not negative after applying promo.
	if totalPrice < 0 {
		return entity.OrderResponse{}, fmt.Errorf("total price cannot be negative after applying promo")
	}

	payload.TotalPrice = totalPrice
	payload.OrderStatus = "preparing"
	payload.Date = time.Now()

	// Validate the fields provided in the payload
	if err := payload.Validate(); err != nil{
		return entity.OrderResponse{}, err
	}

	// Get the customer's balance
	balance, err := uc.balanceRepo.GetUserBalance(payload.CustomerId)
	if err != nil || balance < 0 {
		return entity.OrderResponse{}, fmt.Errorf("failed to get balance")
	}

	// Ensure the total price is less than customer's balance
	if payload.TotalPrice > balance{
		return entity.OrderResponse{}, fmt.Errorf("insufficient balance to complete order")
	}

	// Build a description of ordered items, joining each with commas and "and" for the last item.
	var description string
	for i, item := range payload.OrderItems{
		itemDescription := fmt.Sprintf("%d %s", item.Quantity, item.MenuName)

		if description != "" {
			if i == len(payload.OrderItems)-1 {
					description += " and " + itemDescription 
			} else {
					description += ", " + itemDescription 
			}
	} else {
			description = itemDescription
	}
	}

	// Set up the balance deduction entry with details of the transaction.
	balancePayload := entity.Balance{
		CustomerId:      payload.CustomerId,
		TransactionType: "debit",
		Amount:          payload.TotalPrice,
		Description: "buy " + description,
	}
	
	// Adjust customer's balance by subtracting the order total.
	balancePayload.Balance = balance - balancePayload.Amount
	

	// Insert the value from balancePayload into balances
	_, err = uc.balanceRepo.CreateBalance(balancePayload)
	if err != nil{
		return entity.OrderResponse{}, fmt.Errorf("failed to buy this order: %v", err.Error())
	}

	// Insert the value into orders
	order, err := uc.repo.CreateOrder(payload)
	if err != nil {
			return entity.OrderResponse{}, fmt.Errorf("failed to create order: %v", err)
	}

	order.Date = time.Now().Format("January 02, 2006 03:04 PM")


	// Update promo_used to True 
	if payload.PromoCode != "" {
		err := uc.promoRepo.MarkPromoAsUsed(order.Id)
		if err != nil {
			return entity.OrderResponse{}, fmt.Errorf("failed to mark promo as used: %v", err)
		}
	}

	return order, nil
}

func (uc *orderUseCase) GetUnfinishCustomerOrder(customerId string) (entity.OrderResponse, error){
	return uc.repo.GetUnfinishOrderbyCustomerId(customerId)
}

func (uc *orderUseCase) UpdateOrderStatus(payload entity.OrderResponse) (entity.OrderResponse, error){
	// Retrieve the current order by id
	order, err := uc.repo.GetOrderById(payload.Id)
	if err != nil{
		return entity.OrderResponse{}, err
	}

	// Update the order status to the next stage if applicable
	if order.OrderStatus == "preparing"{
		order.OrderStatus = "out for delivery"
	} else if order.OrderStatus == "out for delivery"{
		order.OrderStatus = "delivered"
	} else{
		return entity.OrderResponse{}, fmt.Errorf("can't update order")
	}

	return uc.repo.UpdateOrderStatus(order)
}

func (uc *orderUseCase) CalculateTotalPrice(payload entity.Order) (float64, error) {
	var totalPrice float64 = 0

	// Iterate the order_items
	for _, item := range payload.OrderItems {
			// Retrieve menu details by name
			menu, err := uc.menuRepo.GetMenubyName(item.MenuName)
			if err != nil {
					return 0, fmt.Errorf("failed to retrieve menu details for item %s: %v", item.MenuName, err)
			}

			// Ensure price and quantity are valid
			if menu.Price == 0 {
					return 0, fmt.Errorf("menu with id %s has invalid price", item.MenuName)
			}
			if item.Quantity <= 0 {
					return 0, fmt.Errorf("invalid quantity for menu item %s", item.MenuName)
			}

			// Calculate item total
			itemTotal := menu.Price * float64(item.Quantity)
			totalPrice += itemTotal
	}

	return totalPrice, nil
}

func (uc *orderUseCase) ApplyPromo(payload entity.Order) (entity.Promo, error) {
	// Retrieve the current promo by promo_code
	promo, err := uc.promoRepo.GetPromoByPromoCode(payload.PromoCode)
	if err != nil {
		return entity.Promo{}, err
	}

	// Check if the promo is currently active based on start and end dates.
	if promo.StartDate.After(time.Now()) || promo.EndDate.Before(time.Now()) {
		return entity.Promo{}, fmt.Errorf("promo code %s is not valid at this time", payload.PromoCode)
	}

	// Verify if the promo has already been used by this customer.
	used, err := uc.promoRepo.IsPromoUsed(payload.CustomerId, payload.PromoCode)
	if err != nil {
		return entity.Promo{}, fmt.Errorf("failed to check promo usage usecase: %v", err)
	}
	if used {
		return entity.Promo{}, fmt.Errorf("promo code %s has already been used by this customer", payload.PromoCode)
	}

	return promo, nil
}

func (uc *orderUseCase) GetAllOrder(page, size int, status string) ([]entity.OrderResponse, model.Paging, error) {
	return uc.repo.GetAllOrder(page, size, status)
}

func (uc *orderUseCase) GetOrderHistory(page, size int, startDate, endDate string, customerId string) ([]entity.OrderResponse, model.Paging, error){
	return uc.repo.GetOrderHistory(page, size, startDate, endDate, customerId)
}


func NewOrderUseCase(repo repository.OrderRepository, menuRepo repository.MenuRepository, balanceRepo repository.BalanceRepository, promoRepo repository.PromoRepository) OrderUseCase{
	return &orderUseCase{repo: repo, menuRepo: menuRepo, balanceRepo: balanceRepo, promoRepo: promoRepo}
}