package model

import "github.com/golang-jwt/jwt/v5"

type MyCustomClaims struct {
    UserId string `json:"id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}
