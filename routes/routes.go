package routes

import (
	"restaurant-api/handlers"
	"restaurant-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.POST("/auth/login", handlers.Login)
	r.POST("/auth/register", handlers.Register)
	auth := r.Group("/")

	auth.Use(middleware.AuthMiddleware())
	auth.GET("/restaurants/:res_id/reports/booking/summary", handlers.GetBookingSummary)
}
