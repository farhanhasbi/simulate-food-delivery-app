package usecase

import (
	"fmt"
	"food-delivery-apps/entity"
	"food-delivery-apps/repository"
	"food-delivery-apps/shared/model"
	"time"
)

type reviewUseCase struct{
	reviewRepo repository.ReviewRepository
	orderRepo repository.OrderRepository
}

type ReviewUseCase interface{
	AddReview(payload entity.Review) (entity.ReviewResponse, error)
	GetReview(page, size int) ([]entity.ReviewResponse, model.Paging, error)
	UpdateReview(payload entity.Review) (entity.ReviewResponse, error)
	DeleteReview(id, customerId string) error
}

func (uc *reviewUseCase) AddReview(payload entity.Review) (entity.ReviewResponse, error){
	payload.UpdatedAt = time.Now()

	// Validate the fields provided in the payload
	if err := payload.Validate(); err != nil{
		return entity.ReviewResponse{}, err
	}

	// Get customer id 
	customerId, err := uc.orderRepo.GetCustomerId(payload.OrderId)
	if err != nil{
		return entity.ReviewResponse{}, err
	}

	payload.BuyDate = customerId.CreatedAt

	// Validate the order is belongs to right customer
	if payload.CustomerId != customerId.CustomerId{
		return entity.ReviewResponse{}, fmt.Errorf("the customer doesn't belong to this order")
	}

	// Check if the menu name, customer ID, and order ID match a completed order.
	var count int
	err = uc.orderRepo.CountfinishOrder(payload.MenuName, payload.CustomerId, payload.OrderId, &count)
	if err != nil{
		return entity.ReviewResponse{}, err
	}
	if count <= 0 {
		return entity.ReviewResponse{}, fmt.Errorf("cannot leave a review for an item that was not ordered")
	}

	return uc.reviewRepo.CreateReview(payload)
}

func (uc *reviewUseCase) GetReview(page, size int) ([]entity.ReviewResponse, model.Paging, error){
	return uc.reviewRepo.GetReview(page, size)
}

func (uc *reviewUseCase) UpdateReview(payload entity.Review) (entity.ReviewResponse, error) {
	// Retrieve the current review by id
	review, err := uc.reviewRepo.GetReviewById(payload.Id)
	if err != nil {
		return entity.ReviewResponse{}, err
	}

	// Ensure customer owns the review
	if review.CustomerId != payload.CustomerId {
		return entity.ReviewResponse{}, fmt.Errorf("unauthorized: customer does not own this review")
	}

	// Validate the fields provided in the payload
	if err := payload.ValidateUpdate(); err != nil {
		return entity.ReviewResponse{}, err
	}

	// Only update fields that are present in the payload
	if payload.Rating != 0 {
		review.Rating = payload.Rating
	}
	if payload.Comment != "" {
		review.Comment = payload.Comment
	}
	
	review.UpdatedAt = time.Now().Format("January 02, 2006 03:04 PM")

	return uc.reviewRepo.UpdateReview(review)
}

func (uc *reviewUseCase) DeleteReview(id, customerId string) error{
	// Retrieve the current review by id
	review, err := uc.reviewRepo.GetReviewById(id)
	if err != nil{
		return err
	}

	if review.CustomerId != customerId{
		return fmt.Errorf("unauthorized: customer does not own this review")
	}

	return uc.reviewRepo.DeleteReview(id, customerId)
}

func NewReviewUseCase(reviewRepo repository.ReviewRepository, orderRepo repository.OrderRepository) ReviewUseCase{
	return &reviewUseCase{reviewRepo: reviewRepo, orderRepo: orderRepo}
}
