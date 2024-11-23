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

type balanceRepository struct {
	db *sql.DB
}

type BalanceRepository interface{
	CreateBalance(payload entity.Balance) (entity.BalanceResponse, error)
	GetBalanceData(page, size int, customerId string) ([]entity.BalanceResponse, model.Paging, error)
	GetUserBalance(customerId string) (float64, error)
}

func (r *balanceRepository) CreateBalance(payload entity.Balance) (entity.BalanceResponse, error){
	// Insert the value for balances
	if err := r.db.QueryRow(config.CreateBalanceQuery, payload.CustomerId, payload.TransactionType,
		payload.Amount, payload.Description, payload.Balance).Scan(&payload.Id, &payload.Balance, &payload.CreatedAt); err != nil{
		return entity.BalanceResponse{}, fmt.Errorf("failed to create balance: %v", err.Error())
	}
	
	// Format CreatedAt for the response in a readable format.
	formattedCreatedAt := payload.CreatedAt.Format("January 02, 2006 03:04 PM")

	// Construct the response object with formatted data.
	response := entity.BalanceResponse{
		Id: payload.Id,
		CustomerId: payload.CustomerId,
		TransactionType: payload.TransactionType,
		Amount: payload.Amount,
		Description: payload.Description,
		Balance: payload.Balance,
		CreatedAt: formattedCreatedAt,
	}

	return response, nil
}

func (r *balanceRepository) GetBalanceData(page, size int, customerId string) ([]entity.BalanceResponse, model.Paging, error){
	var balances []entity.BalanceResponse

	// Calculate the offset for pagination based on the current page and page size.
	offset := (page - 1) * size

	// retrieve the balances with pagination filter by customerId
	rows, err := r.db.Query(config.GetAllUserBalanceQuery, size, offset, customerId)
	if err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to retrieve balance: %v", err.Error())
	}
	defer rows.Close()

	// Iterate over the rows from the database, scanning each row into a balance object.
	for rows.Next(){
		var balance entity.BalanceResponse
		var createdAt time.Time
	
		// Scan balance data into struct fields, including timestamps for creation.
		if err := rows.Scan(&balance.Id, &balance.CustomerId,
			&balance.TransactionType, &balance.Amount, &balance.Description, &balance.Balance, &createdAt); err != nil{
				return nil, model.Paging{}, fmt.Errorf("failed to scan balance: %v", err.Error())
			}

		// Format the CreatedAt for the response in a readable format.
		balance.CreatedAt = createdAt.Format("January 02, 2006 03:04 PM")

		// Append the balance object to the balances slice.
		balances = append(balances, balance)
	}
	
	// Count the total number of balances to set up paging information.
	totalRows := 0
	if err := r.db.QueryRow(config.CountUserBalanceQuery, customerId).Scan(&totalRows); err != nil{
		return nil, model.Paging{}, fmt.Errorf("failed to count user's balance: %v", err.Error())
	}

	// Construct the paging object based on the total rows, page, and size.
	paging := model.Paging{
		Page: page,
		RowsPerPage: size,
		TotalRows: totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
	}

	return balances, paging, nil
}

func (r *balanceRepository) GetUserBalance(customerId string) (float64, error){
	var balance float64

	// Retrieve user's balance by customerId
	err := r.db.QueryRow(config.GetUserBalanceQuery, customerId).Scan(&balance)

	// Handle potential errors from the query
	if err != nil{
		// If no rows are found, return a specific "user not found" error message
		if err == sql.ErrNoRows{
			return 0, fmt.Errorf("balance with customerId %s is not found: %v", customerId, err.Error())
		}
		// For other errors, return a general retrieval failure message
		return 0, fmt.Errorf("failed to retrieve user balance; %v", err.Error())
	}

	return balance, nil
}

func NewBalanceRepository(db *sql.DB) BalanceRepository{
	return &balanceRepository{db: db}
}