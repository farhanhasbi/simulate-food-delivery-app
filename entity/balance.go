package entity

import (
	"fmt"
	"food-delivery-apps/config"
	"time"
)

type Balance struct{
	Id string `json:"id"`
	CustomerId string `json:"-"`
	TransactionType string `json:"transaction_type"`
	Amount float64 `json:"amount"`
	Description string `json:"description"`
	Balance float64 `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type BalanceResponse struct{
	Id string `json:"id"`
	CustomerId string `json:"-"`
	TransactionType string `json:"transaction_type"`
	Amount float64 `json:"amount"`
	Description string `json:"description"`
	Balance float64 `json:"balance"`
	CreatedAt string `json:"created_at"`
}

func (b *Balance) Validate() error{
	if b.Amount == 0 || b.Description == ""{
		return config.ErrMissingFields
	}

	if b.TransactionType != "debit" && b.TransactionType != "credit"{
		return config.ErrInvalidTransactionType
	}

	if b.Amount != 0{
		if b.Amount < 0{
			return fmt.Errorf("amount cannot be below zero")
		}
		if b.Amount < 1000{
			return fmt.Errorf("minimum amount is thousand")
		}
	}
	
	return nil
}