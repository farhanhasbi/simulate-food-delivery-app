package delivery

import (
	"database/sql"
	"fmt"
	"food-delivery-apps/config"
	"food-delivery-apps/usecase"

	"food-delivery-apps/delivery/controller"
	"food-delivery-apps/delivery/middleware"
	"food-delivery-apps/delivery/schedule"
	"food-delivery-apps/repository"
	"food-delivery-apps/shared/service"

	docs "food-delivery-apps/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine *gin.Engine
	host string
	userUc usecase.UserUseCase
	authUc usecase.AuthUseCase
	menuUc usecase.MenuUseCase
	orderUc usecase.OrderUseCase
	balanceUc usecase.BalanceUseCase
	reviewUc usecase.ReviewUseCase
	promoUc usecase.PromoUseCase
	jwtService service.JwtService
}

func (s *Server) initRoute(){	
	docs.SwaggerInfo.BasePath = "/api/v1"
	rg := s.engine.Group(config.ApiGroup)

	// Public Routes
	controller.NewAuthController(s.authUc, rg).Route()
	controller.NewPublicController(s.menuUc, s.reviewUc, rg).Route()

	// Admin Routes
	adminRg := s.engine.Group(config.ApiGroup)
	adminRg.Use(middleware.JWTAuthMiddlewareWithRole(s.jwtService, s.userUc, []string{"admin"}))
	controller.NewAdminController(s.userUc, adminRg).Route()

	// Authentication Routes
	userRg := s.engine.Group(config.ApiGroup)
	userRg.Use(middleware.JWTAuthMiddlewareWithRole(s.jwtService, s.userUc, []string{"admin", "customer", "employee"}))
	controller.NewUserController(s.userUc, userRg).Route()

	// Employee Routes
	employeeRg := s.engine.Group(config.ApiGroup)
	employeeRg.Use(middleware.JWTAuthMiddlewareWithRole(s.jwtService, s.userUc, []string{"employee"}))
	controller.NewEmployeeController(s.menuUc, s.orderUc, s.promoUc, employeeRg).Route()

	// Customer Routes
	customerRg := s.engine.Group(config.ApiGroup)
	customerRg.Use(middleware.JWTAuthMiddlewareWithRole(s.jwtService, s.userUc, []string{"customer"}))
	controller.NewCustomerController(s.orderUc, s.balanceUc, s.reviewUc, s.promoUc, customerRg).Route()

	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *Server) Run(){
	// Initialize routes for the server
	s.initRoute()

	// Start the server and handle any errors
	if err := s.engine.Run(s.host); err != nil{
		panic(fmt.Errorf("error in server %s: %v", s.host, err.Error()))
	}
}

func NewServer() *Server{
	// Load configuration from the config files
	cfg, _ := config.NewConfig()

	// Build the database connection string using the configuration
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	// Open a connection to the database
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil{
		panic(fmt.Errorf("failed to connect database: %v", err.Error()))
	}

	// Initialize JWT service and repositories for various use cases
	jwtService := service.NewJWTService(cfg.TokenConfig)
	userRepo := repository.NewUserRepository(db)
	userUc := usecase.NewUserUseCase(userRepo)
	authUc := usecase.NewAuthUseCase(userUc, jwtService)

	menuRepo := repository.NewMenuRepository(db)
	menuUc := usecase.NewMenuUseCase(menuRepo)

	balanceRepo := repository.NewBalanceRepository(db)
	balanceUc := usecase.NewBalanceUseCase(balanceRepo)

	promoRepo := repository.NewPromoRepository(db)
	promoUc := usecase.NewPromoUseCase(promoRepo)

	orderRepo := repository.NewOrderRepository(db)
	orderUc := usecase.NewOrderUseCase(orderRepo, menuRepo, balanceRepo, promoRepo)

	reviewRepo := repository.NewReviewRepository(db)
	reviewUc := usecase.NewReviewUseCase(reviewRepo, orderRepo)

	// Set up the Gin engine for routing and middleware
	engine := gin.Default()
	
	// Start a background job for periodic tasks
	go schedule.StartCronJob(userUc)
	
	// Define the host and return the server instance with all initialized components
	host := fmt.Sprintf(":%s", cfg.Apiport)
	return &Server{
		engine: engine,
		host: host,
		userUc: userUc,
		authUc: authUc,
		menuUc: menuUc,
		orderUc: orderUc,
		balanceUc: balanceUc,
		reviewUc: reviewUc,
		promoUc: promoUc,
		jwtService: jwtService,
	}
}
