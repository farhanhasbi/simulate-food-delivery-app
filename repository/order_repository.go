package repository

import (
	"database/sql"
	"fmt"
	"food-delivery-apps/config"
	"food-delivery-apps/entity"
	"food-delivery-apps/shared/model"
	"math"
	"time"

	"github.com/lib/pq"
)

type orderRepository struct{
	db *sql.DB
}

type OrderRepository interface{
	CreateOrder(payload entity.Order) (entity.OrderResponse, error)
	CountUnfishOrder(customerId string, count *int) error
	GetUnfinishOrderbyCustomerId(customerId string) (entity.OrderResponse, error)
	GetOrderById(id string) (entity.OrderResponse, error)
	UpdateOrderStatus(payload entity.OrderResponse) (entity.OrderResponse, error) 
	CountfinishOrder(menuId, customerId, orderId string, count *int) error
	GetCustomerId(id string) (entity.Order, error)
	GetAllOrder(page, size int, status string) ([]entity.OrderResponse, model.Paging, error)
	GetOrderHistory(page, size int, startDate, endDate string, customerId string) ([]entity.OrderResponse, model.Paging, error)
}

func (r *orderRepository) CreateOrder(payload entity.Order) (entity.OrderResponse, error){
	// Begin a new transaction.
	tx, err := r.db.Begin()
	if err != nil{
		return entity.OrderResponse{}, fmt.Errorf("failed to begin transaction")
	}

	defer func ()  {
		// Roll back on panic or error, otherwise commit the transaction.
		if p := recover(); p != nil{
			tx.Rollback()
			panic(p)
		} else if err != nil{
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Insert the value for order
	if err := tx.QueryRow(config.CreateOrderQuery, payload.CustomerId, payload.Address, payload.PromoCode, payload.OrderStatus,
		payload.Note, payload.Date, payload.TotalPrice).Scan(&payload.Id, &payload.CreatedAt); err != nil{
			return entity.OrderResponse{}, fmt.Errorf("failed to create order: %v", err.Error())
		}
	
	// Retrieve the customer's username based on CustomerId.
	var username string
	usernameQuery := "SELECT username FROM users WHERE id = $1" 
	if err := r.db.QueryRow(usernameQuery, payload.CustomerId).Scan(&username); err != nil{
		return entity.OrderResponse{}, fmt.Errorf("failed to retrieve username: %v", err.Error())
	}
	payload.CustomerId = username

	// Insert each order item
	for i := range payload.OrderItems{
		payload.OrderItems[i].OrderId = payload.Id

		// Retrieve the menu ID based on the provided menu name.
		var menuId string
		GetMenuNameQuery := "SELECT id FROM menus WHERE name = $1"
		err := r.db.QueryRow(GetMenuNameQuery, payload.OrderItems[i].MenuName).Scan(&menuId)
		if err != nil || menuId == ""{
			return entity.OrderResponse{}, fmt.Errorf("menu with name %s is not found: %v", payload.OrderItems[i].MenuName, err.Error())
		}
		
		payload.OrderItems[i].MenuName = menuId

		// Insert the value for order_items 
		if err := tx.QueryRow(config.CreateOrderItemQuery, payload.OrderItems[i].OrderId,
			payload.OrderItems[i].MenuName, payload.OrderItems[i].Quantity).Scan(&payload.OrderItems[i].Id); err != nil{
				return entity.OrderResponse{}, fmt.Errorf("failed to create order items: %v", err.Error())
			}
	}

	// Convert menu Ids back to names for the response
	for i := range payload.OrderItems {
		var menuName string
		menuQuery := "SELECT name FROM menus WHERE id = $1"
		if err := r.db.QueryRow(menuQuery, payload.OrderItems[i].MenuName).Scan(&menuName); err != nil {
			return entity.OrderResponse{}, fmt.Errorf("failed to retrieve menu name: %v", err.Error())
		}
		payload.OrderItems[i].MenuName = menuName
	}

	// Format timestamps for the response in a readable format.
	formattedCreatedAt := payload.CreatedAt.Format("January 02, 2006 03:04 PM")
	formattedDate := payload.Date.Format("January 02, 2006 03:04 PM")

	// Construct the response object with formatted data.
	response := entity.OrderResponse{
		Id: payload.Id,
		CustomerName: payload.CustomerId,
		Address: payload.Address,
		PromoCode: payload.PromoCode,
		OrderStatus: payload.OrderStatus,
		Note: payload.Note,
		Date: formattedDate,
		TotalPrice: payload.TotalPrice,
		CreatedAt: formattedCreatedAt,
		OrderItems: payload.OrderItems,
	}

	return response, nil
}

func (r *orderRepository) CountUnfishOrder(customerId string, count *int) error{
	if err := r.db.QueryRow(config.CountunfinishCustomerOrderQuery, customerId).Scan(count); err != nil{
		return fmt.Errorf("failed to count unfinished order: %v", err)
	}

	return nil
}

func (r *orderRepository) GetUnfinishOrderbyCustomerId(customerId string) (entity.OrderResponse, error){
	var order entity.OrderResponse

	// retrieve unfinish order by customerId
	err := r.db.QueryRow(config.GetUnfinishOrderByCustomerIdQuery, customerId).Scan(&order.Id, &order.Address,
		&order.PromoCode, &order.OrderStatus, &order.Note, &order.TotalPrice, &order.CreatedAt)
		
		// Handle potential errors from the query 
		if err != nil{
			// If no rows are found, return a specific "customer not has unfinis order" error message
			if err == sql.ErrNoRows{
				return entity.OrderResponse{}, fmt.Errorf("customer with id %s not has unfinish order: %v", customerId, err.Error())
			}
			// For other errors, return a general retrieval failure message
			return entity.OrderResponse{}, fmt.Errorf("failed to retrieve order: %v", err.Error())
		}

	// Retrieve order_items by order_id
	rows, err := r.db.Query(config.GetOrderItemsByOrderIdQuery, order.Id)
	if err != nil{
		return entity.OrderResponse{}, fmt.Errorf("failed to retrieve order items: %v", err.Error())
	}
	defer rows.Close()

	order.OrderItems = []entity.OrderItem{}

	// Iterate over the rows from the database, scanning each row into a orderItem object.
	for rows.Next(){
		var orderItem entity.OrderItem

		// Scan orderItem data into struct fields
		if err := rows.Scan(&orderItem.Id, &orderItem.OrderId, &orderItem.MenuName, &orderItem.Quantity); err != nil {
			return entity.OrderResponse{}, fmt.Errorf("failed to scan order item: %v", err.Error())
		}

		// Append the orderItem object to the orderItems slice.
		order.OrderItems = append(order.OrderItems, orderItem)
	}

	// Parse and format the createdAt for the response in a readable format.
	parsedTime, err := time.Parse(time.RFC3339, order.CreatedAt)
	if err != nil {
		return entity.OrderResponse{}, fmt.Errorf("error parsing time: %v", err.Error())
	}
	formattedCreatedAt := parsedTime.Format("January 02, 2006 03:04 PM")

	// Construct the response object with formatted data.
	response := entity.OrderResponse{
		Id: order.Id,
		Address: order.Address,
		PromoCode: order.PromoCode,
		OrderStatus: order.OrderStatus,
		Note:  order.Note,
		TotalPrice: order.TotalPrice,
		CreatedAt: formattedCreatedAt,
		OrderItems: order.OrderItems,
	}

	return response, nil
}

func (r *orderRepository) GetOrderById(id string) (entity.OrderResponse, error){
	var order entity.OrderResponse

	err := r.db.QueryRow(config.GetOrderByIdQuery, id).Scan(&order.Id, &order.CustomerName, &order.Address,
		&order.PromoCode, &order.OrderStatus, &order.Note, &order.TotalPrice, &order.CreatedAt)
	if err != nil{
		if err == sql.ErrNoRows{
			return entity.OrderResponse{}, fmt.Errorf("order with id %s is not found: %v", id, err.Error())
		}
		return entity.OrderResponse{}, fmt.Errorf("failed to retrieve order: %v", err.Error())
	}

	rows, err := r.db.Query(config.GetOrderItemsByOrderIdQuery, order.Id)
	if err != nil{
		return entity.OrderResponse{}, fmt.Errorf("failed to retrieve order items: %v", err.Error())
	}
	defer rows.Close()

	order.OrderItems = []entity.OrderItem{}
	for rows.Next(){
		var orderItem entity.OrderItem

		if err := rows.Scan(&orderItem.Id, &orderItem.OrderId,
			&orderItem.MenuName, &orderItem.Quantity); err != nil{
				return entity.OrderResponse{}, fmt.Errorf("faild to scan order items: %v", err.Error())
		}

		order.OrderItems = append(order.OrderItems, orderItem)
	}

	if err := rows.Err(); err != nil {
		return entity.OrderResponse{}, fmt.Errorf("error occurred during row iteration: %v", err.Error())
	}

	parsedTime, err := time.Parse(time.RFC3339, order.CreatedAt)
	if err != nil {
		return entity.OrderResponse{}, fmt.Errorf("error parsing time: %v", err.Error())
	}

	formattedCreatedAt := parsedTime.Format("January 02, 2006 03:04 PM")

	response := entity.OrderResponse{
		Id: order.Id,
		CustomerName: order.CustomerName,
		Address: order.Address,
		PromoCode: order.PromoCode,
		OrderStatus: order.OrderStatus,
		Note:  order.Note,
		TotalPrice: order.TotalPrice,
		CreatedAt: formattedCreatedAt,
		OrderItems: order.OrderItems,
	}

	return response, nil
}

func (r *orderRepository) UpdateOrderStatus(payload entity.OrderResponse) (entity.OrderResponse, error){
	_, err := r.db.Exec(config.UpdateOrderStatusQuery, payload.Id, payload.OrderStatus)
	if err != nil{
		return entity.OrderResponse{}, fmt.Errorf("failed to update order status: %v", err.Error())
	}

	return payload, nil
}

func (r *orderRepository) CountfinishOrder(menuId, customerId, orderId string, count *int) error{
	if err := r.db.QueryRow(config.CountfinishCustomerOrderQuery, menuId, customerId, orderId).Scan(count); err != nil{
		return fmt.Errorf("failed to count finish order: %v", err.Error())
	}

	return nil
}

func (r *orderRepository) GetCustomerId(id string) (entity.Order, error){
	var order entity.Order

	// Retrieve customer_id by id and order_status = 'delivered'
	err := r.db.QueryRow(config.GetCustomerIdWithFinishOrderQuery, id).Scan(&order.CustomerId, &order.CreatedAt)

	// Handle potential errors from the query
	if err != nil{
		// If no rows are found, return a specific "order not found" error message
		if err == sql.ErrNoRows{
			return entity.Order{}, fmt.Errorf("order with id %s is not found: %v", id, err.Error())
		}
		// For other errors, return a general retrieval failure message
		return entity.Order{}, fmt.Errorf("faild to find customer id: %v", err.Error())
	}

	return order, nil
}

func (r *orderRepository) GetAllOrder(page, size int, status string) ([]entity.OrderResponse, model.Paging, error){
	var orders []entity.OrderResponse

	// Calculate the offset for pagination based on the current page and page size.
	offset := (page - 1) *size
	
	var statuses []string
	if status == "unfinish"{
		statuses = []string{"preparing", "out for delivery"}
	} else if status == "finish"{
		statuses = []string{"delivered"}
	} else if status == "all" {
		statuses = []string{"preparing", "out for delivery", "delivered"}
	} else {
		return nil, model.Paging{}, fmt.Errorf("status '%s' is not supported for filtering", status)
	}

	// Retrieve all unfinish order with pagination
	rows, err := r.db.Query(config.GetAllOrderQuery, size, offset, pq.Array(statuses))
	
	if err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to retrieve all unfinish order: %v", err.Error())
	}
	defer rows.Close()

	// Iterate over the rows from the database, scanning each row into a order object.
	for rows.Next(){
		var order entity.OrderResponse
		var createdAt time.Time

		// Scan order data into struct fields, including timestamps for creation.
		if err := rows.Scan(&order.Id, &order.CustomerName, &order.Address, &order.PromoCode, &order.OrderStatus, &order.Note, &order.TotalPrice, &createdAt); err != nil{
			return nil, model.Paging{}, fmt.Errorf("failed to scan order: %v", err.Error())
		}

		// Format the timestamps for the response in a readable format.
		order.CreatedAt = createdAt.Format("January 02, 2006 03:04 PM")

		// Retrieve order_items by order_id
		detailrows, err := r.db.Query(config.GetOrderItemsByOrderIdQuery, order.Id)
		if err != nil{
			return nil, model.Paging{}, fmt.Errorf("failed to retrieve order items: %v", err.Error())
		}
		defer detailrows.Close()
		orderItems := []entity.OrderItem{}

		// Iterate over the rows from the database, scanning each row into a orderItem object.
		for detailrows.Next(){
			var orderItem entity.OrderItem

			// Scan orderItem data into struct fields.
			if err := detailrows.Scan(&orderItem.Id, &orderItem.OrderId,
				&orderItem.MenuName, &orderItem.Quantity); err != nil{
					return nil, model.Paging{}, fmt.Errorf("failed to scan order items: %v", err.Error())
				}

			// Append the orderItem object to the orderItems slice.
			orderItems = append(orderItems, orderItem)
		}

		// Assign orderItems into order.OrderItems
		order.OrderItems = orderItems

		// Append the order object to the orders slice.
		orders = append(orders, order)
	}

	// Count the total number of unfinish order to set up paging information.
	totalRows := 0
	if err := r.db.QueryRow(config.CountAllOrderQuery).Scan(&totalRows); err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to count order: %v", err.Error())
	}

	// Construct the paging object based on the total rows, page, and size.
	paging := model.Paging{
		Page: page,
		RowsPerPage: size,
		TotalRows: totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return orders, paging, nil
}

func (r *orderRepository) GetOrderHistory(page, size int, startDate, endDate string, customerId string) ([]entity.OrderResponse, model.Paging, error){
	var orders []entity.OrderResponse

	// Calculate the offset for pagination based on the current page and page size.
	offset := (page - 1) *size
	var rows *sql.Rows
	var err error

	// Retrieve finish order with pagination, otherwise include filter by startDate and endDate
	if startDate != "" && endDate != ""{
		rows, err = r.db.Query(config.GetFilterDateCustomerOrderHistoryQuery, size, offset, startDate, endDate, customerId)
	} else {
		rows, err = r.db.Query(config.GetCustomerOrderHistoryQuery, size, offset, customerId)
	}
	
	if err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to retrieve order history: %v", err.Error())
	}
	defer rows.Close()

	// Iterate over the rows from the database, scanning each row into a order object.
	for rows.Next(){
		var order entity.OrderResponse
		var createdAt, date time.Time

		// Scan order data into struct fields, including timestamps for creation and date for filter purpose.
		if err := rows.Scan(&order.Id, &order.Address, &order.PromoCode, &order.OrderStatus, &order.Note, &date, &order.TotalPrice, &createdAt); err != nil{
			return nil, model.Paging{}, fmt.Errorf("failed to scan order history: %v", err.Error())
		}

		// Format the timestamps for the response in a readable format.
		order.CreatedAt = createdAt.Format("January 02, 2006 03:04 PM")
		order.Date = date.Format("January 02, 2006")

		// Retrieve order_items by order_id
		detailRows, err := r.db.Query(config.GetOrderItemsByOrderIdQuery, order.Id)
		if err != nil{
			return nil, model.Paging{}, fmt.Errorf("failed to retrive order item history: %v", err.Error())
		}
		defer detailRows.Close()
		orderItems := []entity.OrderItem{}

		// Iterate over the rows from the database, scanning each row into a orderItem object.
		for detailRows.Next(){
			var orderItem entity.OrderItem

			// Scan orderItem data into struct fields.
			if err := detailRows.Scan(&orderItem.Id, &orderItem.OrderId,
				&orderItem.MenuName, &orderItem.Quantity); err != nil{
					return nil, model.Paging{}, fmt.Errorf("failed to scan order items: %v", err.Error())
				}

			// Append the orderItem object to the orderItems slice.
			orderItems = append(orderItems, orderItem)
		}

		// Assign orderItems into order.OrderItems
		order.OrderItems = orderItems

		// Append the order object to the orders slice.
		orders = append(orders, order)
	}

	// Count the total number of unfinish order to set up paging information.
	totalRows := 0
	if err := r.db.QueryRow(config.CountFinishCustomerOrderQuery, customerId).Scan(&totalRows); err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to count order history: %v", err.Error())
	}

	// Construct the paging object based on the total rows, page, and size.
	paging := model.Paging{
		Page: page,
		RowsPerPage: size,
		TotalRows: totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return orders, paging, nil
}


func NewOrderRepository(db *sql.DB) OrderRepository{
	return &orderRepository{db: db}
}
