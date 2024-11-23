package model

import (
	"food-delivery-apps/entity"
)

type PromoRequest struct{
	PromoCode string `json:"promo_code"`
	Discount float64 `json:"discount"`
	IsPercentage bool `json:"is_percentage"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Description string `json:"description,omitempty"`
}

type SinglePromoResponse struct{
	Status Status `json:"status"`
	Data entity.PromoResponse `json:"data"`
}

type PagedPromoResponse struct{
	Status Status `json:"status"`
	Data entity.PromoResponse `json:"data"`
	Paging Paging `json:"paging"`
}