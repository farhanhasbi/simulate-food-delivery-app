package entity

import (
	"fmt"
	"food-delivery-apps/config"
	"time"
)

type Menu struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Desc string `json:"description"`
	UnitType string `json:"unit_type"`
	Price float64 `json:"price"`
	Rating float64 `json:"rating"`
	CreatedBy string `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MenuResponse struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Desc string `json:"description"`
	UnitType string `json:"unit_type"`
	Price float64 `json:"price"`
	Rating float64 `json:"rating"`
	CreatedBy string `json:"-"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (m *Menu) Validate() error{
	if m.Name == "" || m.Type == "" || m.Desc == "" || m.UnitType == "" || m.Price == 0{
		return config.ErrMissingFields
	}

	if m.Type != "main dish" && m.Type != "side dish" && m.Type != "dessert" && m.Type != "beverage"{
		return config.ErrInvalidMenuType
	}

	if m.UnitType != "piece" && m.UnitType != "portion" && m.UnitType != "packet" && m.UnitType != "cup"{
		return config.ErrInvalidUnitType
	}

	if m.Price != 0{
		if m.Price < 0{
			return fmt.Errorf("price cannot be below zero")
		}
		if m.Price < 1000{
			return fmt.Errorf("minimum price is 1000")
		}
	}
	
	return nil
}

func (m *Menu) ValidateUpdate() error{
	
	if m.Type != ""{
		if m.Type != "main dish" && m.Type != "side dish" && m.Type != "dessert" && m.Type != "beverage"{
			return config.ErrInvalidMenuType
		}
	}

	if m.UnitType != ""{
		if m.UnitType != "piece" && m.UnitType != "portion" && m.UnitType != "packet" && m.UnitType != "cup"{
			return config.ErrInvalidUnitType
		}
	}
	
	if m.Price != 0{
		if m.Price < 0{
			return fmt.Errorf("price cannot be below zero")
		}
		if m.Price < 500{
			return fmt.Errorf("minimum price is 500")
		}
	}
	
	return nil
}