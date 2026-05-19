package routes

import (
	"restaurant-api/handlers"
	"restaurant-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.POST("/login", handlers.Login)

	auth := r.Group("/")

	auth.Use(middleware.AuthMiddleware())

<<<<<<< Updated upstream
=======
	auth.GET("/booking/:user_id", handlers.GetBookings)

>>>>>>> Stashed changes
}
