package model

import (
	"food-delivery-apps/entity"
)

type CreateReviewRequest struct{
	MenuName string `json:"menu_name"`
	OrderId string 	`json:"order_id"`
	Rating int `json:"rating"`
	Comment string `json:"comment"`
}

type UpdateReviewRequest struct{
	Rating int `json:"rating"`
	Comment string `json:"comment"`
}

type SingleReviewResponse struct{
	Status Status `json:"status"`
	Data entity.ReviewResponse `json:"data"`
}

type PagedReviewResponse struct{
	Status Status `json:"status"`
	Data entity.ReviewResponse `json:"data"`
	Paging Paging `json:"paging"`
}