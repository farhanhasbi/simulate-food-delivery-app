package usecase

import (
	"fmt"
	"food-delivery-apps/entity"
	"food-delivery-apps/repository"
	"food-delivery-apps/shared/model"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	repo repository.UserRepository
}

type UserUseCase interface {
	CreateNewUser(payload entity.User) (entity.UserResponse, error)
	FindUserByEmailPassword(email, password string) (entity.User, error)
	GetAllUser(page, size int, role string) ([]entity.GetUserResponse, model.Paging, error)
	AssignToEmployee(id string) (entity.UserResponse, error)
	UpdateUser(payload entity.User) (entity.UserResponse, error)
	DeleteUser(id string) error
	Logout(token string) error
	IsTokenBlacklisted(token string) bool
	CleanUpExpiredTokens() (int64, error)
}

func (uc *userUseCase) CreateNewUser(payload entity.User) (entity.UserResponse, error){
	// Check if the email already used
	emailExist, _ := uc.repo.GetUserbyEmail(payload.Email)
	if emailExist.Email == payload.Email {
		return entity.UserResponse{}, fmt.Errorf("user with email: %s already exists", payload.Email)
	}

	// Check if the username already used
	userExist, _ := uc.repo.GetUserbyUsername(payload.Username)
	if userExist.Username == payload.Username {
		return entity.UserResponse{}, fmt.Errorf("username: %s already exists", payload.Username)
	}

	// Assign user's role to admin whenever the count equal to 0
	var userCount int
	err := uc.repo.CountUser(&userCount)
	if err != nil {
			return entity.UserResponse{}, err
	}
	if userCount == 0 {
		payload.Role = "admin"
	} else {
		payload.Role = "customer"
	}
	
	payload.UpdatedAt = time.Now()

	// Hashing the password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return entity.UserResponse{}, fmt.Errorf("failed to hash password: %v", err.Error())
	}
	payload.Password = string(hashPassword)

	return uc.repo.CreateNewUser(payload)
}

func (uc *userUseCase) FindUserByEmailPassword(email, password string) (entity.User, error) {
	// Check if user is exist
	userExist, err := uc.repo.GetUserbyEmail(email)
	if err != nil {
		return entity.User{}, fmt.Errorf("user doesn't exists")
	}

	// compare the password from payload and userExist
	err = bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(password))
	if err != nil {
		return entity.User{}, fmt.Errorf("password is Invalid")
	}

	return userExist, nil
}

func (uc *userUseCase) GetAllUser(page, size int, role string) ([]entity.GetUserResponse, model.Paging, error){
	return uc.repo.GetAllUser(page, size, role)
}

func (uc *userUseCase) AssignToEmployee(id string) (entity.UserResponse, error){
	// Retrieve the current user by id
	user, err := uc.repo.GetUserbyId(id)
	if err != nil {
		return entity.UserResponse{}, err
	}

	// assign user's role into given role
	if user.Role == "admin"{
		return entity.UserResponse{}, fmt.Errorf("admins can't be update to employee")
	}
	if user.Role == "employee"{
		return entity.UserResponse{}, fmt.Errorf("user role already update to employee")
	}
	user.Role = "employee"

	user.UpdatedAt = time.Now().Format("January 02, 2006 03:04 PM")

	return uc.repo.UpdateRole(user)
}

func (uc *userUseCase) UpdateUser(payload entity.User) (entity.UserResponse, error){
	// Retrieve the current user by id
	user, err := uc.repo.GetUserbyId(payload.Id)
	if err != nil {
		return entity.UserResponse{}, err
	}

	// Validate the fields provided in the payload
	if err := payload.ValidateUpdate(); err != nil{
		return entity.UserResponse{}, err
	}

	// Check if email is already used
	if payload.Email != "" && payload.Email != user.Email {
		emailExist, _ := uc.repo.GetUserbyEmail(payload.Email)
		if emailExist.Email == payload.Email {
			return entity.UserResponse{}, fmt.Errorf("user with email %s already exists", payload.Email)
		}
	}

	// Check if username is already used
	if payload.Username != "" && payload.Username != user.Username {
		userExist, _ := uc.repo.GetUserbyUsername(payload.Username)
		if userExist.Username == payload.Username {
			return entity.UserResponse{}, fmt.Errorf("username %s already exists", payload.Username)
		}
	}

	// Check if fields are present before updating them
	if payload.Username != "" {
		user.Username = payload.Username
	}
	if payload.Email != "" {
		user.Email = payload.Email
	}
	if payload.Gender != "" {
		user.Gender = payload.Gender
	} 
	
	user.UpdatedAt = time.Now().Format("January 02, 2006 03:04 PM")

	return uc.repo.UpdateUser(user)
}

func (uc *userUseCase) DeleteUser(id string) error{
	// retrieve the current user by id
	_, err := uc.repo.GetUserbyId(id)
	if err != nil{
		return err
	}

	return uc.repo.DeleteUser(id)
}

func (uc *userUseCase) Logout(token string) error {
	// Check if the token is already blacklisted
	if uc.IsTokenBlacklisted(token) {
		return fmt.Errorf("token is already blacklisted")
	}

	// Blacklist the token
	if err := uc.repo.BlackListToken(token); err != nil {
		return fmt.Errorf("failed to logout: %v", err)
	}

	return nil
}

func (uc *userUseCase) IsTokenBlacklisted(token string) bool{
	return uc.repo.IsTokenBlacklisted(token)
}

func (uc *userUseCase) CleanUpExpiredTokens() (int64, error){
	return uc.repo.CleanUpExpiredTokens()
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}
