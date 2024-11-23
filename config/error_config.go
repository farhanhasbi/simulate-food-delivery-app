package config

import (
	"errors"
)

var (
	ErrUserExists      = errors.New("user with email already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidEmail = errors.New("invalid email")
	ErrMissingFields   = errors.New("some required fields are missing")
	ErrInvalidGender   = errors.New("gender must be either male or female")
	ErrInvalidMenuType = errors.New("menu type must be main dish, side dish, dessert or beverage")
	ErrInvalidRole = errors.New("role must be admin, employee or customer")
	ErrInvalidOrderStatus = errors.New("order status must be preparing, out for delivery, or delivered")
	ErrInvalidTransactionType = errors.New("transaction type be either debit or credit")
	ErrInvalidUnitType = errors.New("unit type must be piece, portion, packet or cup")
)