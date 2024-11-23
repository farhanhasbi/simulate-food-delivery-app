package entity

import (
	"fmt"
	"food-delivery-apps/config"
	"time"
)

type Review struct{
	Id string `json:"id"`
	CustomerId string `json:"customer_id"`
	MenuName string `json:"menu_name"`
	OrderId string 	`json:"order_id"`
	Rating int `json:"rating"`
	Comment string `json:"comment"`
	BuyDate time.Time `json:"buy_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReviewRequest struct{
	MenuName string `json:"menu_name"`
	OrderId string 	`json:"order_id"`
	Rating int `json:"rating"`
	Comment string `json:"comment"`
}

type ReviewResponse struct{
	Id string `json:"id"`
	CustomerId string `json:"-"`
	CustomerName string `json:"customer_name"`
	MenuName string `json:"menu_name"`
	OrderId string 	`json:"order_id"`
	Rating int `json:"rating"`
	Comment string `json:"comment"`
	BuyDate string `json:"buy_date"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (r *Review) Validate() error{
	if r.MenuName == "" || r.OrderId == "" || r.Comment == ""{
		return config.ErrMissingFields
	}

	if r.Rating == 0 {
		if r.Rating < 0 {
			return fmt.Errorf("rating cannot be below zero")
		}
		if r.Rating > 5 {
			return fmt.Errorf("rating cannot exceed 5")
		}
		return config.ErrMissingFields
	}

	return nil
}

func (r *Review) ValidateUpdate() error {
	if r.Rating != 0 {
		if r.Rating < 0 {
			return fmt.Errorf("rating cannot be below zero")
		}
		if r.Rating > 5 {
			return fmt.Errorf("rating cannot exceed 5")
		}
	}
	
	return nil
}
