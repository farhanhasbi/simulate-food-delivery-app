package repository

import (
	"database/sql"
	"fmt"
	"food-delivery-apps/config"
	"food-delivery-apps/entity"
	"food-delivery-apps/shared/model"
	"math"
	"time"
)

type reviewRepository struct{
	db *sql.DB
}

type ReviewRepository interface{
	CreateReview(payload entity.Review) (entity.ReviewResponse, error)
	GetReview(page, size int) ([]entity.ReviewResponse, model.Paging, error)
	GetReviewById(id string) (entity.ReviewResponse, error)
	UpdateReview(payload entity.ReviewResponse) (entity.ReviewResponse, error)
	DeleteReview(id, customerId string) error
}

func (r *reviewRepository) CreateReview(payload entity.Review) (entity.ReviewResponse, error){
	// Retrieve the menu ID based on the provided menu name.
	var menuId string
	getMenuIdQuery := "SELECT id FROM menus WHERE name = $1"
	err := r.db.QueryRow(getMenuIdQuery, payload.MenuName).Scan(&menuId)
	if err != nil || menuId == ""{
		return entity.ReviewResponse{}, fmt.Errorf("menu with name %s is not found: %v", payload.MenuName, err.Error())
	}
	payload.MenuName = menuId
	
	// Insert the value for review.
	if err := r.db.QueryRow(config.CreateReviewQuery, payload.CustomerId,
		payload.MenuName, payload.OrderId, payload.Rating,
		payload.Comment, payload.BuyDate, payload.UpdatedAt).Scan(&payload.Id, &payload.CreatedAt);	err != nil{
		return entity.ReviewResponse{}, fmt.Errorf("failed to create review: %v", err.Error())
	}

	// Retrieve the customer's username based on CustomerId.
	var customerName string
	getCustomerNameQuery := "SELECT username FROM users WHERE id = $1"
	if err := r.db.QueryRow(getCustomerNameQuery, payload.CustomerId).Scan(&customerName); err != nil{
		return entity.ReviewResponse{}, fmt.Errorf("failed to get customer username: %v", err.Error())
	}
	payload.CustomerId = customerName

	// Convert menu IDs back to names for the response
	var menuName string
	getMenuNameQuery := "SELECT name FROM menus WHERE id = $1"
	if err := r.db.QueryRow(getMenuNameQuery, payload.MenuName).Scan(&menuName); err != nil{
		return entity.ReviewResponse{}, fmt.Errorf("failed to get menu name: %v", err.Error())
	}
	payload.MenuName = menuName
	
	// Format timestamps for the response in a readable format.
	formattedCreatedAt := payload.CreatedAt.Format("January 02, 2006 03:04 PM")
	formattedUpdatedAt := payload.UpdatedAt.Format("January 02, 2006 03:04 PM")
	formattedBuyDate := payload.BuyDate.Format("January 02, 2006 03:04 PM")

	// Construct the response object with formatted data.
	response := entity.ReviewResponse{
		Id: payload.Id,
		CustomerName: payload.CustomerId,
		MenuName: payload.MenuName,
		OrderId: payload.OrderId,
		Rating: payload.Rating,
		Comment: payload.Comment,
		BuyDate: formattedBuyDate,
		CreatedAt: formattedCreatedAt,
		UpdatedAt: formattedUpdatedAt,
	}

	return response, nil
}

func (r *reviewRepository) GetReview(page, size int) ([]entity.ReviewResponse, model.Paging, error){
	var reviews []entity.ReviewResponse

	// Calculate the offset for pagination based on the current page and page size.
	offset := (page - 1) *size

	// Retrieve review with pagination
	rows, err := r.db.Query(config.GetAllReviewQuery, size, offset)
	if err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to retrieve reviews: %v", err.Error())
	}
	defer rows.Close()

	// Iterate over the rows from the database, scanning each row into a review object.
	for rows.Next(){
		var review entity.ReviewResponse
		var buyDate, createdAt, updatedAt time.Time

		// Scan review data into struct fields, including timestamps for creation, update, and buyDate.
		if err := rows.Scan(&review.Id, &review.CustomerName, &review.MenuName,
			&review.OrderId, &review.Rating, &review.Comment, &buyDate, &createdAt, &updatedAt); err != nil{
				return nil, model.Paging{}, fmt.Errorf("failed to scan reviews: %v", err.Error())
			}

		// Format the timestamps for the response in a readable format.
		review.BuyDate = buyDate.Format("January 02, 2006 03:04 PM")
		review.CreatedAt = createdAt.Format("January 02, 2006 03:04 PM")
		review.UpdatedAt = updatedAt.Format("January 02, 2006 03:04 PM")

		// Append the review object to the reviews slice.
		reviews = append(reviews, review)
	}

	// Count the total number of reviews to set up paging information.
	totalRows := 0
	if err := r.db.QueryRow(config.CountReviewQuery).Scan(&totalRows); err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to count review: %v", err.Error())
	}

	// Construct the paging object based on the total rows, page, and size.
	paging := model.Paging{
		Page: page,
		RowsPerPage: size,
		TotalRows: totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return reviews, paging, nil
}

func (r *reviewRepository) GetReviewById(id string) (entity.ReviewResponse, error){
	var review entity.ReviewResponse
	var customerId string

	// Retrieve review by id
	err := r.db.QueryRow(config.GetReviewByIdQuery, id).Scan(&review.Id, &customerId, &review.CustomerName,
		&review.MenuName, &review.OrderId, &review.Rating,
		&review.Comment, &review.BuyDate, &review.CreatedAt, &review.UpdatedAt)

	// Handle potential errors from the query
	if err != nil{
		// If no rows are found, return a specific "menu not found" error message
		if err == sql.ErrNoRows{
			return entity.ReviewResponse{}, fmt.Errorf("review with id %s is not found: %v", id, err.Error())
		}
		// For other errors, return a general retrieval failure message
		return entity.ReviewResponse{}, fmt.Errorf("failed to retrieve review: %v", err.Error())
	}

	review.CustomerId = customerId

	// Parse and format the CreatedAt for the response in a readable format.
	parsedCreatedAt, err := time.Parse(time.RFC3339, review.CreatedAt)
	if err != nil {
		return entity.ReviewResponse{}, fmt.Errorf("error parsing time: %v", err.Error())
	}
	formattedCreatedAt := parsedCreatedAt.Format("January 02, 2006 03:04 PM")

	// Parse and format the BuyDate for the response in a readable format.
	parsedBuyDate, err := time.Parse(time.RFC3339, review.BuyDate)
	if err != nil {
		return entity.ReviewResponse{}, fmt.Errorf("error parsing time: %v", err.Error())
	}
	formattedBuyDate := parsedBuyDate.Format("January 02, 2006 03:04 PM")

	// Parse and format the UpdatedAt for the response in a readable format.
	parsedUpdatedAt, err := time.Parse(time.RFC3339, review.UpdatedAt)
	if err != nil {
		return entity.ReviewResponse{}, fmt.Errorf("error parsing time: %v", err.Error())
	}
	formattedUpdatedAt := parsedUpdatedAt.Format("January 02, 2006 03:04 PM")

	// Construct the response object with formatted data.
	response := entity.ReviewResponse{
		Id: review.Id,
		CustomerId: review.CustomerId,
		CustomerName: review.CustomerName,
		MenuName: review.MenuName,
		OrderId: review.OrderId,
		Rating: review.Rating,
		Comment: review.Comment,
		BuyDate: formattedBuyDate,
		CreatedAt: formattedCreatedAt,
		UpdatedAt: formattedUpdatedAt,
	}

	return response, nil
}

func (r *reviewRepository) UpdateReview(payload entity.ReviewResponse) (entity.ReviewResponse, error){
	_, err := r.db.Exec(config.UpdateReviewQuery, payload.Id, payload.CustomerId, payload.Rating, payload.Comment, payload.UpdatedAt)
	if err != nil{
		return entity.ReviewResponse{}, fmt.Errorf("failed to update review: %v", err.Error())
	}

	return payload, nil
}

func (r *reviewRepository) DeleteReview(id, customerId string) error{
	_, err := r.db.Exec(config.DeleteReviewQuery, id, customerId)
	if err != nil{
		return fmt.Errorf("failed to delete review: %v", err.Error())
	}

	return nil
}

func NewReviewRepository(db *sql.DB) ReviewRepository{
	return &reviewRepository{db: db}
}