package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

// AppConfig contains the application configurations
type AppConfig struct {
	DB            *sql.DB
	jwtSecret     []byte
	encryptionKey []byte
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", composePostgresConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize AppConfig with the database connection
	appConfig := &AppConfig{
		DB:            db,
		jwtSecret:     []byte(os.Getenv("JWT_SECRET")),
		encryptionKey: []byte(os.Getenv("ENCRYPTION_KEY")),
	}

	// Initialize Echo router
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// Handlers
	registerHandlers(e, appConfig)

	// Start server
	fmt.Println("Server started on port 8080")
	e.Start(":8080")
}

// registerHandlers registers all HTTP handlers
func registerHandlers(e *echo.Echo, appConfig *AppConfig) {
	// customerRepo := repository.NewCustomerRepository(appConfig.DB)
	// transactionRepo := repository.NewTransactionRepository(appConfig.DB)
	// limitRepo := repository.NewLimitRepository(appConfig.DB)

	// authHandler := handler.NewAuthHandler(customerRepo, appConfig.jwtSecret, appConfig.encryptionKey)
	// transactionHandler := handler.NewTransactionHandler(transactionRepo, limitRepo)
	// limitHandler := handler.NewLimitHandler(limitRepo, customerRepo)

	// e.POST("/auth/register", authHandler.RegisterCustomer)
	// e.POST("/auth/login", authHandler.LoginHandler)

	// fundGroup := e.Group("/fund")
	// fundGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningKey: appConfig.jwtSecret,
	// }))

	// fundGroup.POST("/transaction", transactionHandler.CreateTransaction)
	// fundGroup.POST("/limit", limitHandler.CreateLimit)
}

// composePostgresConnectionString creates the PostgreSQL connection string
func composePostgresConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
}
