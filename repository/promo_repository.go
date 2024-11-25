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

type promoRepository struct{
	db *sql.DB
}

type PromoRepository interface{
	CreatePromo(payload entity.PromoRequest) (entity.PromoResponse, error)
	GetAllPromo(page, size int) ([]entity.PromoResponse, model.Paging, error)
	GetPromoForCustomer(page, size int, customerId string) ([]entity.PromoResponse, model.Paging, error)
	GetPromoByPromoCode(code string) (entity.Promo, error)
	GetPromoById(id string) (entity.PromoResponse, error)
	DeletePromo(id string) error
	IsPromoUsed(customerID, promoCode string) (bool, error)
	MarkPromoAsUsed(orderID string) error
}

func (r *promoRepository) CreatePromo(payload entity.PromoRequest) (entity.PromoResponse, error){
	const layout = "2006-01-02"
	var createdAt, updatedAt time.Time

	// Insert the value for promos
	if err := r.db.QueryRow(config.CreatePromoQuery, payload.EmployeeId, payload.PromoCode,
		payload.Discount, payload.IsPercentage, payload.StartDate, payload.EndDate, payload.Description, payload.UpdatedAt).Scan(&payload.Id, &createdAt, &updatedAt); err != nil{

		if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" { // Unique constraint violation
						return entity.PromoResponse{}, fmt.Errorf("promo with name %s already exists", payload.PromoCode)
				}
				return entity.PromoResponse{}, fmt.Errorf("database error: %s", pqErr.Message)
		}
		
		return entity.PromoResponse{}, fmt.Errorf("failed to create promo: %v", err.Error())
	}
	
	// Parse and format the startDate for the response in a readable format.
	startDate, err := time.Parse(layout, payload.StartDate)
	if err != nil {
		return entity.PromoResponse{}, fmt.Errorf("invalid start date format: %v", err)
	}
	formattedStartDate := startDate.Format("January 02, 2006 03:04 PM")
	
	// Parse and format the endDate for the response in a readable format.
	endDate, err := time.Parse(layout, payload.EndDate)
	if err != nil {
		return entity.PromoResponse{}, fmt.Errorf("invalid end date format: %v", err)
	}
	formattedEndDate := endDate.Format("January 02, 2006 03:04 PM")

	// Format the createdAt and updatedAt for the response in a readable format.
	formattedCreatedAt := createdAt.Format("January 02, 2006 03:04 PM")
	formattedUpdatedAt := updatedAt.Format("January 02, 2006 03:04 PM")

	// Construct the response object with formatted data.
	response := entity.PromoResponse{
		Id:  payload.Id,
		EmployeeId: payload.EmployeeId,
		PromoCode: payload.PromoCode,
		Discount: payload.Discount,
		IsPercentage: payload.IsPercentage,
		StartDate: formattedStartDate,
		EndDate: formattedEndDate,
		Description: payload.Description,
		CreatedAt: formattedCreatedAt,
		UpdatedAt: formattedUpdatedAt,
	}

	return response, nil
}

func (r *promoRepository) GetAllPromo(page, size int) ([]entity.PromoResponse, model.Paging, error){
	var promos []entity.PromoResponse

	// Calculate the offset for pagination based on the current page and page size.
	offset := (page - 1) *size

	// Retrieve promo with pagination
	rows, err := r.db.Query(config.GetAllPromoQuery, size, offset)
	if err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to retrieve promo: %v", err.Error())
	}
	defer rows.Close()

	// Iterate over the rows from the database, scanning each row into a promo object.
	for rows.Next(){
		var promo entity.PromoResponse
		var startDate, endDate, createdAt, updatedAt time.Time

		// Scan promo data into struct fields, including timestamps for creation, update, start and end date.
		if err := rows.Scan(&promo.Id, &promo.PromoCode, &promo.Discount,
			&promo.IsPercentage, &startDate, &endDate, &promo.Description, &createdAt, &updatedAt); err != nil{
				return nil, model.Paging{}, fmt.Errorf("failed to scan promo: %v", err.Error())
			}

		// Format the timestamps for the response in a readable format.
		promo.StartDate = startDate.Format("January 02, 2006 03:04 PM")
		promo.EndDate = endDate.Format("January 02, 2006 03:04 PM")
		promo.CreatedAt = createdAt.Format("January 02, 2006 03:04 PM")
		promo.UpdatedAt = updatedAt.Format("January 02, 2006 03:04 PM")

		// Append the promo object to the promos slice.
		promos = append(promos, promo)
	}

	// Count the total number of promos to set up paging information.
	totalRows := 0
	if err := r.db.QueryRow(config.CountPromoQuery).Scan(&totalRows); err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to count promo: %v", err.Error())
	}

	// Construct the paging object based on the total rows, page, and size.
	paging := model.Paging{
		Page: page,
		RowsPerPage: size,
		TotalRows: totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return promos, paging, nil
}

func (r *promoRepository) GetPromoForCustomer(page, size int, customerId string) ([]entity.PromoResponse, model.Paging, error){
	var promos []entity.PromoResponse

	// Calculate the offset for pagination based on the current page and page size.
	offset := (page - 1) *size

	// Retrieve promo with pagination
	rows, err := r.db.Query(config.GetPromoForCustomerQuery, size, offset, customerId)
	if err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to retrieve promo: %v", err.Error())
	}
	defer rows.Close()

	// Iterate over the rows from the database, scanning each row into a promo object.
	for rows.Next(){
		var promo entity.PromoResponse
		var startDate, endDate, createdAt, updatedAt time.Time

		// Scan promo data into struct fields, including timestamps for creation, update, start and end date.
		if err := rows.Scan(&promo.Id, &promo.PromoCode, &promo.Discount,
			&promo.IsPercentage, &startDate, &endDate, &promo.Description, &createdAt, &updatedAt); err != nil{
				return nil, model.Paging{}, fmt.Errorf("failed to scan promo: %v", err.Error())
			}

		// Format the timestamps for the response in a readable format.
		promo.StartDate = startDate.Format("January 02, 2006 03:04 PM")
		promo.EndDate = endDate.Format("January 02, 2006 03:04 PM")
		promo.CreatedAt = createdAt.Format("January 02, 2006 03:04 PM")
		promo.UpdatedAt = updatedAt.Format("January 02, 2006 03:04 PM")

		// Append the promo object to the promos slice.
		promos = append(promos, promo)
	}

	// Count the total number of promos to set up paging information.
	totalRows := 0
	if err := r.db.QueryRow(config.CountPromoForCustomerQuery, customerId).Scan(&totalRows); err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to count promo: %v", err.Error())
	}

	// Construct the paging object based on the total rows, page, and size.
	paging := model.Paging{
		Page: page,
		RowsPerPage: size,
		TotalRows: totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return promos, paging, nil
}

func (r *promoRepository) GetPromoByPromoCode(code string) (entity.Promo, error) {
	var promoResponse entity.PromoResponse

	// Retrieve promo by promo_code
	err := r.db.QueryRow(config.GetPromobyPromoCodeQuery, code).Scan(&promoResponse.Id, &promoResponse.PromoCode, &promoResponse.Discount,
		&promoResponse.IsPercentage, &promoResponse.StartDate, &promoResponse.EndDate, &promoResponse.Description)

	// Handle potential errors from the query
	if err != nil {
		// If no rows are found, return a specific "promo not found" error message
		if err == sql.ErrNoRows {
			return entity.Promo{}, fmt.Errorf("promo with name %s is not found: %v", code, err.Error())
		}
		// For other errors, return a general retrieval failure message
		return entity.Promo{}, fmt.Errorf("failed to retrieve promo: %v", err.Error())
	}

	// Parse the StartDate for the response in a readable format.
	parsedStartDate, err := time.Parse(time.RFC3339, promoResponse.StartDate)
	if err != nil {
		return entity.Promo{}, fmt.Errorf("error parsing start date: %v", err.Error())
	}

	// Parse the EndDate for the response in a readable format.
	parsedEndDate, err := time.Parse(time.RFC3339, promoResponse.EndDate)
	if err != nil {
		return entity.Promo{}, fmt.Errorf("error parsing end date: %v", err.Error())
	}

	// Construct the response object with formatted data.
	promo := entity.Promo{
		Id:          promoResponse.Id,
		PromoCode:   promoResponse.PromoCode,
		Discount:    promoResponse.Discount,
		IsPercentage: promoResponse.IsPercentage,
		StartDate:   parsedStartDate,
		EndDate:     parsedEndDate,
		Description: promoResponse.Description,
	}

	return promo, nil
}

func (r *promoRepository) GetPromoById(id string) (entity.PromoResponse, error){
	var promo entity.PromoResponse

	// Retrieve promo by id
	err := r.db.QueryRow(config.GetPromoByIdQuery, id).Scan(&promo.Id)
	if err != nil{
		// If no rows are found, return a specific "promo not found" error message
		if err == sql.ErrNoRows{
			return entity.PromoResponse{}, fmt.Errorf("promo with id %s is not found: %v", id, err.Error())
		}
		// For other errors, return a general retrieval failure message
		return entity.PromoResponse{}, fmt.Errorf("failed to retrieve promo: %v", err.Error())
	}

	return promo, nil
}

func (r *promoRepository) DeletePromo(id string) error{
	_, err := r.db.Exec(config.DeletePromoQuery, id)
	if err != nil{
		return fmt.Errorf("failed to delete promo: %v", err.Error())
	}

	return nil
}

func (r *promoRepository) IsPromoUsed(customerID, promoCode string) (bool, error){
	var count int

	if err := r.db.QueryRow(config.CountUsagePromoQuery, customerID, promoCode).Scan(&count); err != nil{
		return false, fmt.Errorf("failed to check promo usage repo: %v", err)
	}

	return count > 0, nil
}

func (r *promoRepository) MarkPromoAsUsed(orderID string) error {
	_, err := r.db.Exec(config.UpdatePromoUsedStatusQuery, orderID)
	return err
}


func NewPromoRepository(db *sql.DB) PromoRepository{
	return &promoRepository{db: db}
}