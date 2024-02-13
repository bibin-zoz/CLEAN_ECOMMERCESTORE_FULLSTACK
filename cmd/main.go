// main.go
package main

import (
	database "cleancode/db"
	"cleancode/delivery/handlers"
	respository "cleancode/respository"
	"cleancode/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	db := database.InitDB()
	// Initialize repository, use case, and handler

	userRepository := respository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepository)
	userHandler := handlers.NewUserHandler(userUseCase)

	// Set up Gin router
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*.html")
	router.Static("/static", "./static")
	// Define API routes
	router.GET("/register", userHandler.Signup)
	router.POST("/register", userHandler.RegisterUser)
	api := router.Group("/api")
	{
		userAPI := api.Group("/user")
		{
			userAPI.POST("/register", userHandler.RegisterUser)

		}
	}

	// Run the application
	router.Run(":8080")
}
