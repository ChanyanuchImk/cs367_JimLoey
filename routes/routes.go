package routes

import (
	"restaurant-api/handlers"
	"restaurant-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.POST("/auth/login", handlers.Login)

	auth := r.Group("/")

	auth.Use(middleware.AuthMiddleware())
	auth.GET("/restaurants/:res_id/reports/booking/summary", handlers.GetBookingSummary)

	auth.POST("/booking", handlers.CreateBooking)

	auth.GET("/booking/:res_id", handlers.GetBookings)
	auth.PATCH("/booking/:res_id/:book_id/status", handlers.UpdateBookingStatus)
	auth.PUT("/booking/:res_id/:book_id/status", handlers.UpdateBookingStatus)
	auth.POST("/restaurants", handlers.CreateRestaurant)
	auth.GET("/bookings/queue", handlers.GetQueues)

	auth.GET("/booking/:user_id", handlers.GetBookings)

}
