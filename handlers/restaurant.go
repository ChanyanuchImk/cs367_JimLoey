package handlers

import (
	"net/http"
	"restaurant-api/database"

	"github.com/gin-gonic/gin"
)

type CreateRestaurantRequest struct {
	OwnerID     int    `json:"owner_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Phone       string `json:"phone"`
	OpenTime    string `json:"open_time"`
	CloseTime   string `json:"close_time"`
}

func CreateRestaurant(c *gin.Context) {
	var req CreateRestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.OwnerID == 0 || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing restaurant information"})
		return
	}

	result, err := database.DB.Exec(
		`
		INSERT INTO RESTAURANTS (owner_id, name, description, location, phone, open_time, close_time)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
		req.OwnerID,
		req.Name,
		req.Description,
		req.Location,
		req.Phone,
		req.OpenTime,
		req.CloseTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
		return
	}

	restaurantID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get restaurant id"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"restaurant_id": restaurantID,
		"owner_id":      req.OwnerID,
		"name":          req.Name,
		"description":   req.Description,
		"location":      req.Location,
		"phone":         req.Phone,
		"open_time":     req.OpenTime,
		"close_time":    req.CloseTime,
		"status":        "active",
	})
}
