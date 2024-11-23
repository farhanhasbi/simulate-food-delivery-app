package usecase

import (
	"fmt"
	"food-delivery-apps/entity"
	"food-delivery-apps/repository"
	"food-delivery-apps/shared/model"
	"time"
)

type promoUseCase struct{
	repo repository.PromoRepository
}

type PromoUseCase interface{
	CreatePromo(payload entity.PromoRequest) (entity.PromoResponse, error)
	GetAllPromo(page, size int) ([]entity.PromoResponse, model.Paging, error)
	GetPromoForCustomer(page, size int, customerId string) ([]entity.PromoResponse, model.Paging, error)
	DeletePromo(id string) error
}

func (uc *promoUseCase) CreatePromo(payload entity.PromoRequest) (entity.PromoResponse, error){
	// convert promoRequest to promo
	promo, err := payload.ToPromo()
	if err != nil {
		return entity.PromoResponse{}, err
	}

	// Validate the fields provided in the payload
	if err := promo.Validate(); err != nil{
		return entity.PromoResponse{}, err
	}

	// Check if the promo_code is already used
	promoExist, _ := uc.repo.GetPromoByPromoCode(payload.PromoCode)
	if promoExist.PromoCode == payload.PromoCode{
		return entity.PromoResponse{}, fmt.Errorf("promo with code %s already exists", payload.PromoCode)
	}

	payload.UpdatedAt = time.Now().Format("January 02, 2006 03:04 PM")

	return uc.repo.CreatePromo(payload)
}

func (uc *promoUseCase) GetAllPromo(page, size int) ([]entity.PromoResponse, model.Paging, error){
	return uc.repo.GetAllPromo(page, size)
}

func (uc *promoUseCase) GetPromoForCustomer(page, size int, customerId string) ([]entity.PromoResponse, model.Paging, error){
	return uc.repo.GetPromoForCustomer(page, size, customerId)
}

func (uc *promoUseCase) DeletePromo(id string) error{
	// Retrieve the current promo by id
	_, err := uc.repo.GetPromoById(id)
	if err != nil{
		return err
	}

	return uc.repo.DeletePromo(id)
}

func NewPromoUseCase(repo repository.PromoRepository) PromoUseCase{
	return &promoUseCase{repo: repo}
}