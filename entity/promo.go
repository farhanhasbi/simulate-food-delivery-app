package entity

import (
	"errors"
	"fmt"
	"food-delivery-apps/config"
	"time"
)

type Promo struct{
	Id string `json:"id"`
	EmployeeId string `json:"-"`
	PromoCode string `json:"promo_code"`
	Discount float64 `json:"discount"`
	IsPercentage bool `json:"is_percentage"`
	StartDate time.Time `json:"start_date"`
	EndDate time.Time `json:"end_date"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PromoRequest struct{
	Id string `json:"id"`
	EmployeeId string `json:"-"`
	PromoCode string `json:"promo_code"`
	Discount float64 `json:"discount"`
	IsPercentage bool `json:"is_percentage"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Description string `json:"description"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PromoResponse struct{
	Id string `json:"id"`
	EmployeeId string `json:"-"`
	PromoCode string `json:"promo_code"`
	Discount float64 `json:"discount"`
	IsPercentage bool `json:"is_percentage"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Description string `json:"description,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (req *PromoRequest) ToPromo() (Promo, error) {
	const layout = "2006-01-02"
	startDate, err := time.Parse(layout, req.StartDate)
	if err != nil {
		return Promo{}, fmt.Errorf("invalid start date format: %v", err)
	}

	endDate, err := time.Parse(layout, req.EndDate)
	if err != nil {
		return Promo{}, fmt.Errorf("invalid end date format: %v", err)
	}

	var createdAt, updatedAt time.Time

	return Promo{
		Id:           req.Id,
		EmployeeId:   req.EmployeeId,
		PromoCode:    req.PromoCode,
		Discount:     req.Discount,
		IsPercentage: req.IsPercentage,
		StartDate:    startDate,
		EndDate:      endDate,
		Description:  req.Description,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (p *Promo) Validate() error{
	if p.PromoCode == "" || p.Discount == 0{
		return config.ErrMissingFields
	}

	if p.StartDate.After(p.EndDate){
		return fmt.Errorf("start date can't pass the end date")
	}

	if p.StartDate.After(time.Now()) {
		return fmt.Errorf("start date cannot be in the future")
	}

	if p.EndDate.Before(time.Now()) {
		return fmt.Errorf("end date can't be in the past")
	}

	if p.Discount != 0{
		if p.Discount < 0{
			return fmt.Errorf("discount cannot be below zero")
		}
		if p.IsPercentage && p.Discount >= 100 {
			return errors.New("percentage discount cannot exceed 100%")
		}
		if !p.IsPercentage && p.Discount < 10000 {
			return errors.New("minimum discount for promo code without percent is 10000")
		}
	}

	return nil
}