// main.go
package main

import (
	database "cleancode/pkg/db"
	"cleancode/pkg/delivery/handlers"
	"cleancode/pkg/delivery/middleware"
	respository "cleancode/pkg/respository"
	"cleancode/pkg/usecase"

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
	router.GET("/register", middleware.IsLogin(), userHandler.Signup)
	router.POST("/register", userHandler.RegisterUser)
	router.GET("/verify", middleware.IsLogin(), userHandler.VerifyHandler)
	router.POST("/verify", userHandler.VerifyPost)
	router.GET("/login", middleware.IsLogin(), userHandler.LoginHandler)
	router.POST("/login", middleware.IsLogin(), userHandler.LoginPost)
	router.GET("/home", middleware.LoginAuth(), userHandler.HomeHandler)

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
