package main

import (
	"log"
	"os"

	"primo-11/internal/handler"
	"primo-11/internal/repository"
	"primo-11/internal/usecase"
	"primo-11/internal/database"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "primo-11/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type CustomValidator struct {
	validator *validator.Validate
}


func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i) 
}

// @title           Product API 
// @version         1.0 
// @description     Product API with PostgreSQL
// @host            localhost:5000
// @BasePath        / 

func main() {
	_ = godotenv.Load()

	db, err := database.NewPostgresDB()
	if err != nil { 
		log.Fatalf("Database connection failed: %v", err)
	}
	productRepo := repository.NewProductRepository(db)  
	productUC := usecase.NewProductUseCase(productRepo)  
	productHandler := handler.NewProductHandler(productUC) 

	e := echo.New()
	e.Use(middleware.Logger())  
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}

	productHandler.RegisterRoutes(e)
	e.GET("/api-docs/*", echoSwagger.WrapHandler)


	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("Server starting on port %s", port)
	log.Printf("API Docs: http://localhost:%s/api-docs/index.html", port) 

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}


