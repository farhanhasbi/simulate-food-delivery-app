package main

import (
	"food-delivery-apps/delivery"
)

// @title           Food Delivery API
// @version         1.0
// @description     This is a Food Delivery API server using Gin and Clean Architecture.

// @contact.name   Farhan Hasbi
// @contact.email  farhanhasbi512@gmail.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.bearer BearerAuth
// @in header
// @name Authorization
func main() {
	delivery.NewServer().Run()
}