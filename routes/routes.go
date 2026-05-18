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

	auth.PATCH("/booking/:res_id/:book_id/status", handlers.UpdateBookingStatus)
	auth.PUT("/booking/:res_id/:book_id/status", handlers.UpdateBookingStatus)
}
