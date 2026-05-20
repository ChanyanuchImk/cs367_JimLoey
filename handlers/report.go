package handlers

import (
	"net/http"
	"restaurant-api/database"

	"github.com/gin-gonic/gin"
)

type BookingSummary struct {
	TotalBookings     int     `json:"total_bookings"`
	CancelledBookings int     `json:"cancelled_bookings"`
	CompletedBookings int     `json:"completed_bookings"`
	TotalRevenue      float64 `json:"total_revenue"`
	StartDate         string  `json:"start_date"`
	EndDate           string  `json:"end_date"`
}

func GetBookingSummary(c *gin.Context) {
	restaurantID := c.Param("res_id")
	startDate := c.Query("start_date") // ?start_date=2026-01-01
	endDate := c.Query("end_date")     // ?end_date=2026-12-31

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start_date and end_date are required",
		})
		return
	}

	query := `
        SELECT
            COUNT(*) AS total_bookings,
            SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END) AS cancelled_bookings,
            SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) AS completed_bookings,
            COALESCE(SUM(CASE WHEN status != 'cancelled' THEN total_price ELSE 0 END), 0) AS total_revenue
        FROM BOOKINGS
        WHERE restaurant_id = ?
          AND booking_date BETWEEN ? AND ?
    `

	var summary BookingSummary
	err := database.DB.QueryRow(query, restaurantID, startDate, endDate).Scan(
		&summary.TotalBookings,
		&summary.CancelledBookings,
		&summary.CompletedBookings,
		&summary.TotalRevenue,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch summary",
		})
		return
	}

	summary.StartDate = startDate
	summary.EndDate = endDate

	c.JSON(http.StatusOK, summary)
}
