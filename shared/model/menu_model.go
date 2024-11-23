package model

import (
	"food-delivery-apps/entity"
)

type MenuRequest struct{
	Name string `json:"name"`
	Type string `json:"type"`
	Desc string `json:"description"`
	UnitType string `json:"unit_type"`
	Price float64 `json:"price"`
}

type SingleMenuResponse struct{
	Status Status `json:"status"`
	Data entity.MenuResponse `json:"data"`
}

type PagedMenuResponse struct{
	Status Status `json:"status"`
	Data entity.MenuResponse `json:"data"`
	Paging Paging `json:"paging"`
}