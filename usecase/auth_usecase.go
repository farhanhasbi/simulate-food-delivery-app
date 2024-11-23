package usecase

import (
	"fmt"
	"food-delivery-apps/entity"
	"food-delivery-apps/entity/dto"
	"food-delivery-apps/shared/service"
)

type authUseCase struct{
	uc UserUseCase
	jwtservice service.JwtService
}

type AuthUseCase interface{
	Register(payload dto.AuthRequestRegister) (entity.UserResponse, error)
	Login(payload dto.AuthRequestLogin) (dto.AuthResponse, error)
}

func (a *authUseCase) Login(payload dto.AuthRequestLogin) (dto.AuthResponse, error) {
	// Check if the user is exist by email and password
	user, err := a.uc.FindUserByEmailPassword(payload.Email, payload.Password)
	if err != nil {
		return dto.AuthResponse{}, fmt.Errorf("failed to login: %v", err.Error())
	}

	// Create token after login
	token, err := a.jwtservice.CreateToken(user)
	if err != nil {
		return dto.AuthResponse{}, fmt.Errorf("failed to create token: %v", err.Error())
	}

	return token, nil
}

func (a *authUseCase) Register(payload dto.AuthRequestRegister) (entity.UserResponse, error) {
	// Validate the fields provided in the payload
	if err := payload.Validate(); err != nil{
		return entity.UserResponse{}, err
	}
		
	return a.uc.CreateNewUser(entity.User{
		Username: payload.Username,
		Email: payload.Email,
		Password: payload.Password,
		Gender: payload.Gender,
	})
}

func NewAuthUseCase(uc UserUseCase, jwtservice service.JwtService) AuthUseCase {
	return &authUseCase{uc: uc, jwtservice: jwtservice,}
}

/*
"email": "admin@mail.com",
"password": "strongPassword123"
_______________________________
"email": "Yejina25@mail.com",
"password": "psychicLover25"
_______________________________
"email": "shirooni12@mail.com",
"password": "breakthelimit"
_______________________________
"email": "whitelotus@mail.com",
"password": "teenager123"
_______________________________
"email": "antimonitor@mail.com",
"password": "offtherecordnight"
*/