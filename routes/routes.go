package routes

import (
	"restaurant-api/handlers"
	"restaurant-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// นำ ErrorHandler มาดักจับทุกเส้นทาง
	r.Use(middleware.ErrorHandler())

	r.POST("/login", handlers.Login)

	auth := r.Group("/")

	auth.Use(middleware.AuthMiddleware())

}
