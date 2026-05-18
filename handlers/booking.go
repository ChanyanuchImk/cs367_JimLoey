package handlers

import (
	"net/http"
	"restaurant-api/database"

	"github.com/gin-gonic/gin"
)

type Booking struct {
	BookingID      int     `json:"booking_id"`
	UserID         int     `json:"user_id"`
	RestaurantID   int     `json:"restaurant_id"`
	TableID        *int    `json:"table_id"`
	BookingDate    string  `json:"booking_date"`
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	NumberOfPeople int     `json:"number_of_people"`
	TotalPrice     float64 `json:"total_price"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at"`
}

// GET /booking/{res_id}
func GetBookings(c *gin.Context) {
	resID := c.Param("res_id")

	query := `
		SELECT booking_id, user_id, restaurant_id, table_id, booking_date, start_time, end_time, number_of_people, total_price, status, created_at
		FROM BOOKINGS
		WHERE restaurant_id = ?
	`

	rows, err := database.DB.Query(query, resID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query bookings"})
		return
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var b Booking
		var tableID *int
		err := rows.Scan(
			&b.BookingID, &b.UserID, &b.RestaurantID, &tableID,
			&b.BookingDate, &b.StartTime, &b.EndTime,
			&b.NumberOfPeople, &b.TotalPrice, &b.Status, &b.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse booking data"})
			return
		}
		b.TableID = tableID
		bookings = append(bookings, b)
	}

	if bookings == nil {
		bookings = []Booking{}
	}

	c.JSON(http.StatusOK, bookings)
}
