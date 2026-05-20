package main

import (
	"restaurant-api/database"
	"restaurant-api/routes"
	"restaurant-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	database.ConnectDB()

	r := gin.Default()

	r.Use(middleware.ErrorHandler())

	routes.SetupRoutes(r)

	r.Run(":8080")

}
