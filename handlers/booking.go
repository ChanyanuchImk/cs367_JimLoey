package handlers

import (
	"net/http"
	"restaurant-api/database"

	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateBookingRequest struct {
	UserID         int    `json:"user_id"`
	RestaurantID   int    `json:"restaurant_id"`
	BookingDate    string `json:"booking_date"`
	StartTime      string `json:"start_time"`
	NumberOfPeople int    `json:"number_of_people"`
}

// POST /booking
func CreateBooking(c *gin.Context) {
	var req CreateBookingRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.UserID == 0 || req.RestaurantID == 0 || req.BookingDate == "" || req.StartTime == "" || req.NumberOfPeople <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing booking information"})
		return
	}

	startTime, err := time.Parse("15:04:05", req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
		return
	}

	endTime := startTime.Add(2 * time.Hour).Format("15:04:05")

	var restaurantID int
	err = database.DB.QueryRow(
		"SELECT restaurant_id FROM RESTAURANTS WHERE restaurant_id = ?",
		req.RestaurantID,
	).Scan(&restaurantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	var tableID int
	query := `
		SELECT table_id
		FROM TABLES
		WHERE restaurant_id = ?
			AND capacity >= ?
			AND table_id NOT IN (
				SELECT table_id
				FROM BOOKINGS
				WHERE restaurant_id = ?
					AND booking_date = ?
					AND start_time = ?
					AND status != 'cancelled'
					AND table_id IS NOT NULL
			)
		ORDER BY capacity ASC
		LIMIT 1
	`

	err = database.DB.QueryRow(
		query,
		req.RestaurantID,
		req.NumberOfPeople,
		req.RestaurantID,
		req.BookingDate,
		req.StartTime,
	).Scan(&tableID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusConflict, gin.H{"error": "No available table"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find available table"})
		return
	}

	result, err := database.DB.Exec(
		`
		INSERT INTO BOOKINGS (user_id, restaurant_id, table_id, booking_date, start_time, end_time, number_of_people)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
		req.UserID,
		req.RestaurantID,
		tableID,
		req.BookingDate,
		req.StartTime,
		endTime,
		req.NumberOfPeople,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	bookingID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get booking id"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"booking_id":       bookingID,
		"user_id":          req.UserID,
		"restaurant_id":    req.RestaurantID,
		"table_id":         tableID,
		"booking_date":     req.BookingDate,
		"start_time":       req.StartTime,
		"end_time":         endTime,
		"number_of_people": req.NumberOfPeople,
		"status":           "pending",
	})
}

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
		// Scan ข้อมูลให้ครบตามจำนวน Column ใน SELECT
		err := rows.Scan(
			&b.BookingID, &b.UserID, &b.RestaurantID, &b.TableID,
			&b.BookingDate, &b.StartTime, &b.EndTime,
			&b.NumberOfPeople, &b.TotalPrice, &b.Status, &b.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}
		bookings = append(bookings, b)
	}

	if bookings == nil {
		bookings = []Booking{}
	}
	c.JSON(http.StatusOK, bookings)
	// ----------------------------------
}

func GetBookingsByUser(c *gin.Context) {
	userID := c.Param("user_id")
	query := `SELECT booking_id, user_id, restaurant_id, table_id, booking_date, start_time, end_time, number_of_people, total_price, status, created_at FROM BOOKINGS WHERE user_id = ?`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query bookings"})
		return
	}
	defer rows.Close()
	var bookings []Booking
	for rows.Next() {
		var b Booking
		err := rows.Scan(&b.BookingID, &b.UserID, &b.RestaurantID, &b.TableID, &b.BookingDate, &b.StartTime, &b.EndTime, &b.NumberOfPeople, &b.TotalPrice, &b.Status, &b.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse booking data"})
			return
		}
		bookings = append(bookings, b)
	}
	if bookings == nil {
		bookings = []Booking{}
	}
	c.JSON(http.StatusOK, bookings)
}

func GetBookingsByRestaurant(c *gin.Context) {
	resID := c.Param("res_id")
	query := `SELECT booking_id, user_id, restaurant_id, table_id, booking_date, start_time, end_time, number_of_people, total_price, status, created_at FROM BOOKINGS WHERE restaurant_id = ?`
	rows, err := database.DB.Query(query, resID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query bookings"})
		return
	}
	defer rows.Close()
	var bookings []Booking
	for rows.Next() {
		var b Booking
		err := rows.Scan(&b.BookingID, &b.UserID, &b.RestaurantID, &b.TableID, &b.BookingDate, &b.StartTime, &b.EndTime, &b.NumberOfPeople, &b.TotalPrice, &b.Status, &b.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse booking data"})
			return
		}
		bookings = append(bookings, b)
	}
	if bookings == nil {
		bookings = []Booking{}
	}
	c.JSON(http.StatusOK, bookings)
}
