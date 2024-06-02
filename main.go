package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"referral-system/handler"
	"referral-system/middleware"
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

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	fmt.Println("Database connection successful")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Error JWT Secret required.")
	}

	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	referralRepo := repository.NewReferralLinkRepository(db)
	contributionRepo := repository.NewContributionRepository(db)
	dbTransaction := repository.NewDBTransactionRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUsecase(dbTransaction, userRepo, roleRepo, referralRepo, contributionRepo, jwtSecret)
	referralLinkUsecase := usecase.NewReferralLinkUsecase(dbTransaction, referralRepo)

	// Initialize HTTP handlers
	userHandler := handler.NewUserHandler(userUseCase)
	referralLinkHandler := handler.NewReferralHandler(referralLinkUsecase)

	// Initialize Echo router
	e := echo.New()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// Handlers
	registerHandlers(e, jwtSecret, userHandler, referralLinkHandler)

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
func composePostgresConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
}
