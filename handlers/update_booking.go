package handlers

import (
	"net/http"
	"restaurant-api/database"

	"github.com/gin-gonic/gin"
)

type UpdateBookingStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// PATCH /booking/{res_id}/{book_id}/status
func UpdateBookingStatus(c *gin.Context) {
	resID := c.Param("res_id")
	bookID := c.Param("book_id")

	var req UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body. Status is required."})
		return
	}

	// Basic validation for status ENUM
	validStatuses := map[string]bool{"pending": true, "confirmed": true, "completed": true, "cancelled": true}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	query := `
		UPDATE BOOKINGS
		SET status = ?
		WHERE restaurant_id = ? AND booking_id = ?
	`

	result, err := database.DB.Exec(query, req.Status, resID, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check affected rows"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found or status is already the same"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking status updated successfully"})
}
