package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type DbConfig struct{
	Host string
	Port string
	User string
	Password string
	Name string
	Driver string
}

type ApiConfig struct{
	Apiport string
}

type TokenConfig struct {
	IssuerName       string `json:"IssuerName"`
	JwtSignatureKey   []byte `json:"JwtSignatureKy"`
	JwtSigningMethod *jwt.SigningMethodHMAC
	JwtExpiresTime   time.Duration
}

type Config struct{
	DbConfig
	ApiConfig
	TokenConfig
}

func (c *Config) ReadConfig() error {
	// Load the .env file into environment variables
	err := godotenv.Load()
	if err != nil{
		return fmt.Errorf("file .env can't be load: %v", err.Error())
	}

	// Populate the DbConfig struct fields with database configuration values from the .env file
	c.DbConfig = DbConfig{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name: os.Getenv("DB_NAME"),
		Driver: os.Getenv("DB_DRIVER"),
	}
	
	// Populate the DbConfig struct fields with api configuration values (port) from the .env file
	c.ApiConfig = ApiConfig{
		Apiport: os.Getenv("API_PORT"),
	}

	// Parse token expiration time from the .env file as an integer, with default fallback on error
	tokenExpire, _ := strconv.Atoi(os.Getenv("TOKEN_EXPIRE"))
	// Populate the TokenConfig struct with token-related settings from the .env file
	c.TokenConfig = TokenConfig{
		IssuerName: os.Getenv("TOKEN_ISSUER"),
		JwtSignatureKey: []byte(os.Getenv("TOKEN_SECRET")),
		JwtSigningMethod: jwt.SigningMethodHS256,
		JwtExpiresTime: time.Duration(tokenExpire) * time.Minute,
	}

	// Validate all required config fields, ensuring none are empty or invalid
	if c.Host == "" || c.Port == "" || c.User == "" || c.Password == "" || c.Name == "" || c.Driver == "" ||
	 c.Apiport == "" || c.IssuerName == "" || c.JwtExpiresTime < 0 || len(c.JwtSignatureKey) == 0 {
		return fmt.Errorf("some required config is missing")
	}

	return nil
}

func NewConfig() (*Config, error) {
	// Create a new, empty Config instance
	cfg := &Config{}

	// Attempt to load configuration values into the Config instance by calling ReadConfig
	err := cfg.ReadConfig()
	if err != nil{
		return nil, fmt.Errorf("can't load config: %v", err.Error())
	}

	return cfg, nil
}