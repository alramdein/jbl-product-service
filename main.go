package main

import (
	"database/sql"
	"fmt"
	"log"
	"product-service/config"
	"product-service/handler"
	"product-service/middleware"
	"product-service/repository"
	"product-service/usecase"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.Db.Host, conf.Db.Port, conf.Db.User, conf.Db.Password, conf.Db.Name, conf.Db.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not ping the database: %v", err)
	}

	productRepo := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(middleware.JWTMiddleware)

	apiGroup := e.Group("/api")
	handler.NewProductHandler(apiGroup, productUsecase)

	log.Println("Server running on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
