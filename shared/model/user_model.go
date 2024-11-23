package model

import (
	"food-delivery-apps/entity"
	"food-delivery-apps/entity/dto"
)

type UserRequest struct{
	Username string `json:"username"`
	Email string `json:"email"`
	Gender string `json:"gender"`
}

type LoginResponse struct{
	Status Status `json:"status"`
	Data dto.AuthResponse `json:"data"`
}

type SingleUserResponse struct{
	Status Status `json:"status"`
	Data entity.UserResponse `json:"data"`
}

type PagedUserResponse struct{
	Status Status `json:"status"`
	Data entity.BalanceResponse `json:"data"`
	Paging Paging `json:"paging"`
}