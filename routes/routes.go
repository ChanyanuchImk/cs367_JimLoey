package routes

import (
	"restaurant-api/handlers"
	"restaurant-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.POST("/auth/login", handlers.Login)
	// ตรวจสอบว่าใน handlers มี Register หรือไม่ ถ้าไม่มีบรรทัดนี้จะแดง
	// r.POST("/auth/register", handlers.Register)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.GET("/search", handlers.SearchRestaurants)
	auth.GET("/search/:res_id", handlers.GetRestaurantDetail)

	// รายงานสรุปการจอง
	auth.GET("/restaurants/:res_id/reports/booking/summary", handlers.GetBookingSummary)

	auth.POST("/booking", handlers.CreateBooking)

	// การจัดการโต๊ะและคิว
	auth.PATCH("/tables/:table_id/status", handlers.UpdateTableStatus)
	auth.GET("/tables/count", handlers.GetTableCount)
	auth.GET("/bookings/queue", handlers.GetQueues)

	// การจัดการสถานะการจอง
	auth.PATCH("/booking/:res_id/:book_id/status", handlers.UpdateBookingStatus)
	auth.PUT("/booking/:res_id/:book_id/status", handlers.UpdateBookingStatus)

	// สร้างร้านค้าใหม่
	auth.POST("/restaurants", handlers.CreateRestaurant)

	// --- แก้ไขจุดที่แดง (ชื่อฟังก์ชันต้องตรงกับใน booking.go) ---

	// เปลี่ยนชื่อ URL ให้สื่อสารชัดเจนเพื่อไม่ให้ Route ทับกัน
	auth.GET("/bookings/restaurant/:res_id", handlers.GetBookingsByRestaurant) //
	auth.GET("/bookings/user/:user_id", handlers.GetBookingsByUser)            //
}
