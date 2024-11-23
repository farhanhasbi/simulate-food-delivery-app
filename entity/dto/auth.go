package dto

import (
	"fmt"
	"food-delivery-apps/config"
	"regexp"
)

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type AuthRequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRequestRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

func (a AuthRequestRegister) IsValidEmail() bool {
	return EmailRegex.MatchString(a.Email)
}

func (a *AuthRequestRegister) ValidatePassword() error {
	if len(a.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	return nil
}

func (a *AuthRequestRegister) Validate() error {
	if a.Username == "" || a.Email == "" || a.Password == "" || a.Gender == "" {
			return config.ErrMissingFields
	}

	if a.Gender != "male" && a.Gender != "female" {
			return config.ErrInvalidGender
	}

	if !a.IsValidEmail() { 
		return config.ErrInvalidEmail
	}

	if err := a.ValidatePassword(); err != nil {
		return err
	}

	return nil
}