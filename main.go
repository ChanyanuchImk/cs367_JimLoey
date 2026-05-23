package main

import (
	"restaurant-api/database"
	"restaurant-api/middleware"
	"restaurant-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	database.ConnectDB()

	r := gin.Default()

	r.Use(middleware.ErrorHandler())

	routes.SetupRoutes(r)

	r.Run(":8080")

}
