package main

import (
	"restaurant-api/database"
	"restaurant-api/routes"

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

	routes.SetupRoutes(r)

	r.Run(":8080")

}
