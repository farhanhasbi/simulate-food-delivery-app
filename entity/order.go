package entity

import (
	"fmt"
	"food-delivery-apps/config"
	"time"
)

type Order struct{
	Id string `json:"id"`
	CustomerId string `json:"customer_id"`
	Address string `json:"address"`
	PromoCode string `json:"promo_code"`
	OrderStatus string `json:"order_status"`
	Note string `json:"note"`
	Date time.Time `json:"date"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	OrderItems  []OrderItem `json:"order_items"`
}

type OrderResponse struct{
	Id string `json:"id"`
	CustomerName string `json:"customer_name,omitempty"`
	PromoCode string `json:"promo_code,omitempty"`
	Address string `json:"address"`
	OrderStatus string `json:"order_status"`
	Note string `json:"note,omitempty"`
	Date string `json:"date,omitempty"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt  string `json:"created_at"`
	OrderItems  []OrderItem `json:"order_items"`
}

type OrderItem struct{
	Id string `json:"id"`
	OrderId string `json:"-"`
	MenuName string `json:"menu_name"`
	Quantity int `json:"quantity"`
}


func (o *Order) Validate() error{
	if o.Address == "" {
		return config.ErrMissingFields
	}
	
	if o.TotalPrice == 0 {
		return fmt.Errorf("failed to calculate total price")
	}

	if o.OrderStatus != "preparing" && o.OrderStatus != "out for delivery" &&  o.OrderStatus != "delivered"{
		return config.ErrInvalidOrderStatus
	}

	for _, item := range o.OrderItems{
		if item.Quantity <= 0 {
			return fmt.Errorf("can't set quantity to zero or below")
		}
	}

	return nil
}