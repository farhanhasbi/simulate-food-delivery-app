package service

import (
	"fmt"
	"food-delivery-apps/config"
	"food-delivery-apps/entity"
	"food-delivery-apps/entity/dto"
	"food-delivery-apps/shared/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct{
	cfg config.TokenConfig
}

type JwtService interface{
	CreateToken(user entity.User) (dto.AuthResponse, error)
	ParseToken(tokenHeader string) (jwt.MapClaims, error)
	ValidateToken(token string) (jwt.MapClaims, error)
}

func (j *jwtService) CreateToken(user entity.User) (dto.AuthResponse, error){
	// Set the token expiration time based on the configured JwtExpiresTime
	expiresAt := time.Now().Add(j.cfg.JwtExpiresTime)

	// Create custom claims including user details and token expiration/issue times
	claims := model.MyCustomClaims{
		UserId: user.Id,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: j.cfg.IssuerName,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create a new token with claims and signing method
	token := jwt.NewWithClaims(j.cfg.JwtSigningMethod, claims)

	// Sign the token using the secret key from the config
	ss, err := token.SignedString(j.cfg.JwtSignatureKey)
	if err != nil{
		return dto.AuthResponse{}, fmt.Errorf("failed to sign token: %w", err)
	}

	// Calculate token expiration in minutes
	expiresIn := int64(j.cfg.JwtExpiresTime.Minutes())

	// Return the generated token and its expiration time
	return dto.AuthResponse{
		Token: ss,
		ExpiresIn: expiresIn,
	}, nil	
}

func (j *jwtService) ParseToken(tokenHeader string) (jwt.MapClaims, error){
	// Parse the JWT token and validate the signing method
	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for verification
		return j.cfg.JwtSignatureKey, nil
	})

	// Return an error if the token parsing fails
	if err != nil{
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if the token is valid
	if !token.Valid{
		return nil, fmt.Errorf("invalid token")
	}

	// Extract the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("oops, failed to parse token claims")
	}
	// Return the parsed claims
	return claims, nil
}
 
func (j *jwtService) ValidateToken(token string) (jwt.MapClaims, error){
	// Parse and validate the token
	claims, err := j.ParseToken(token)
	if err != nil{
		return nil, err
	}
	// Return the claims if the token is valid
	return claims, nil
}

func NewJWTService(cfg config.TokenConfig) JwtService{
	return &jwtService{cfg: cfg}
}