package main

import (
	"database/sql"
	"fmt"
	"log"
	"referral-system/config"
	"referral-system/handler"
	"referral-system/middleware"
	"referral-system/repository"
	"referral-system/usecase"

	_ "referral-system/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Referral System API
// @version 1.0
// @description Server for referral system.
// @termsOfService http://swagger.io/terms/

// @contact.name Alif Ramdani
// @contact.url https://github.com/alramdein
// @contact.email ramdanialif26@gmail.com
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", composePostgresConnectionString(cfg.Db))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	fmt.Println("Database connection successful")

	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	referralRepo := repository.NewReferralLinkRepository(db)
	contributionRepo := repository.NewContributionRepository(db)
	dbTransaction := repository.NewDBTransactionRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUsecase(dbTransaction, userRepo, roleRepo, referralRepo, contributionRepo, cfg.JwtSecret, cfg.ReferralLinkExp)
	referralLinkUsecase := usecase.NewReferralLinkUsecase(dbTransaction, referralRepo)

	// Initialize HTTP handlers
	userHandler := handler.NewUserHandler(userUseCase)
	referralLinkHandler := handler.NewReferralHandler(referralLinkUsecase)

	// Initialize Echo router
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// Serve Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Handlers
	registerHandlers(e, cfg.JwtSecret, userHandler, referralLinkHandler)

	// Start server
	fmt.Println("Server started on port 8080")
	e.Start(":8080")
}

// registerHandlers registers all HTTP handlers
func registerHandlers(e *echo.Echo, jwtSecret string, userHandler *handler.UserHandler, referralLinkHandler *handler.ReferralHandler) {
	e.POST("/register", userHandler.RegisterUserGenerator)
	e.POST("/register/:code", userHandler.RegisterUserContributor)

	e.POST("/login", userHandler.Login)

	e.POST("/referral-link", referralLinkHandler.GenerateReferralLink, middleware.JWTMiddleware(jwtSecret))
}

// composePostgresConnectionString creates the PostgreSQL connection string
func composePostgresConnectionString(cfg config.DbConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName)
}
