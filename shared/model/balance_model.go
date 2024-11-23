package model

import (
	"food-delivery-apps/entity"
)

type BalanceRequest struct{
	Amount float64 `json:"amount"`
	Description string `json:"description"`
}

type SingleBalanceResponse struct{
	Status Status `json:"status"`
	Data entity.BalanceResponse `json:"data"`
}

type PagedBalanceResponse struct{
	Status Status `json:"status"`
	Data entity.BalanceResponse `json:"data"`
	Paging Paging `json:"paging"`
}