package usecase

import (
	"fmt"
	"food-delivery-apps/entity"
	"food-delivery-apps/repository"
	"food-delivery-apps/shared/model"
)

type balanceUseCase struct{
	repo repository.BalanceRepository
}

type BalanceUseCase interface{
	IncreaseBalance(payload entity.Balance) (entity.BalanceResponse, error)
	GetBalanceData(page, size int, customerId string) ([]entity.BalanceResponse, model.Paging, error)
}

func (uc *balanceUseCase) IncreaseBalance(payload entity.Balance) (entity.BalanceResponse, error){
	payload.TransactionType = "credit"

	// Validate the fields provided in the payload
	if err := payload.Validate(); err != nil{
		return entity.BalanceResponse{}, err
	}

	// Get customer's balance
	balance, _ := uc.repo.GetUserBalance(payload.CustomerId)
	if balance < 0 {
		return entity.BalanceResponse{}, fmt.Errorf("failed to get balance")
	}

	// Adjust customer's balance by adding the amount.
	payload.Balance = balance + payload.Amount

	return uc.repo.CreateBalance(payload)
}

func (uc *balanceUseCase) GetBalanceData(page, size int, customerId string) ([]entity.BalanceResponse, model.Paging, error){
	return uc.repo.GetBalanceData(page, size, customerId)
}

func NewBalanceUseCase(repo repository.BalanceRepository) BalanceUseCase{
	return &balanceUseCase{repo: repo}
}