package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"referral-system/handler"
	"referral-system/repository"
	"referral-system/usecase"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

// AppConfig contains the application configurations
// type AppConfig struct {
// 	DB            *sql.DB
// 	jwtSecret     []byte
// 	encryptionKey []byte
// }

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
	// appConfig := &AppConfig{
	// 	DB:            db,
	// 	jwtSecret:     []byte(os.Getenv("JWT_SECRET")),
	// 	encryptionKey: []byte(os.Getenv("ENCRYPTION_KEY")),
	// }

	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	referralRepo := repository.NewReferralLinkRepository(db)
	contributionRepo := repository.NewContributionRepository(db)
	dbTransaction := repository.NewDBTransactionRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUsecase(dbTransaction, userRepo, roleRepo, referralRepo, contributionRepo)

	// Initialize HTTP handlers
	userHandler := handler.NewUserHandler(userUseCase)

	// Initialize Echo router
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// Handlers
	registerHandlers(e, userHandler)

	// Start server
	fmt.Println("Server started on port 8080")
	e.Start(":8080")
}

// registerHandlers registers all HTTP handlers
func registerHandlers(e *echo.Echo, userHandler *handler.UserHandler) {
	e.POST("/register", userHandler.RegisterUserGenerator)
	e.POST("/register/:code", userHandler.RegisterUserContributor)
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
