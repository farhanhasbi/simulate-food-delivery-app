package model

import (
	"food-delivery-apps/entity"
)

type OrderRequest struct{
	Address string `json:"address"`
	PromoCode string `json:"promo_code"`
	Note string `json:"note"`
	OrderItems  []OrderItemRequest `json:"order_items"`
}

type OrderItemRequest struct{
	MenuName string `json:"menu_name"`
	Quantity int `json:"quantity"`
}

type SingleOrderResponse struct{
	Status Status `json:"status"`
	Data entity.OrderResponse `json:"data"`
}

type PagedOrderResponse struct{
	Status Status `json:"status"`
	Data entity.OrderResponse `json:"data"`
	Paging Paging `json:"paging"`
}