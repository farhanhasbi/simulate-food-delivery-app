package entity

import (
	"food-delivery-apps/config"
	"regexp"
	"time"
)

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type User struct{
	Id string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
  Gender string `json:"gender"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserResponse struct{
	Id string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
  Gender string `json:"gender"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type GetUser struct{
	Id string `json:"id"`
	Username string `json:"username"`
	Role string `json:"role"`
  Gender string `json:"gender"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GetUserResponse struct{
	Id string `json:"id"`
	Username string `json:"username"`
	Role string `json:"role"`
  Gender string `json:"gender"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (u User) IsValidEmail() bool {
	return EmailRegex.MatchString(u.Email)
}

func (u *User) ValidateUpdate() error{
	if u.Gender != ""{
		if u.Gender != "male" && u.Gender != "female"{
			return config.ErrInvalidGender
		}
	}

	if u.Email != ""{
		if !u.IsValidEmail() { 
			return config.ErrInvalidEmail
		}
	}

	return nil
}

