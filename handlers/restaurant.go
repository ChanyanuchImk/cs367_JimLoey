package handlers

import (
	"net/http"
	"restaurant-api/database"

	"github.com/gin-gonic/gin"
)

type RestaurantRequest struct {
	OwnerID     int    `json:"owner_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Phone       string `json:"phone"`
	OpenTime    string `json:"open_time"`
	CloseTime   string `json:"close_time"`
}

func CreateRestaurant(c *gin.Context) {

	var req RestaurantRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})

		return
	}

	query := `
	INSERT INTO RESTAURANTS
	(
		owner_id,
		name,
		description,
		location,
		phone,
		open_time,
		close_time
	)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := database.DB.Exec(
		query,
		req.OwnerID,
		req.Name,
		req.Description,
		req.Location,
		req.Phone,
		req.OpenTime,
		req.CloseTime,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Restaurant created",
		"restaurant_id": id,
	})
}
