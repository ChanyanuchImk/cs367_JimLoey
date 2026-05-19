package handlers

import (
	"database/sql"
	"net/http"
	"restaurant-api/database"

	"github.com/gin-gonic/gin"
)

type Restaurant struct {
	RestaurantID int    `json:"restaurant_id"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	Phone        string `json:"phone"`
	OpenTime     string `json:"open_time"`
	CloseTime    string `json:"close_time"`
}

func SearchRestaurants(c *gin.Context) {

	keyword := c.Query("keyword")
	query := `
		SELECT 
			restaurant_id,
			name,
			location,
			phone,
			open_time,
			close_time
		FROM restaurants
		WHERE status = 'active'
	`

	var rows *sql.Rows
	var err error

	if keyword != "" {
		query += " AND name LIKE ?"
		rows, err = database.DB.Query(query, "%"+keyword+"%")
	} else {
		rows, err = database.DB.Query(query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var restaurants []Restaurant
	for rows.Next() {
		var restaurant Restaurant
		rows.Scan(
			&restaurant.RestaurantID,
			&restaurant.Name,
			&restaurant.Location,
			&restaurant.Phone,
			&restaurant.OpenTime,
			&restaurant.CloseTime,
		)
		restaurants = append(restaurants, restaurant)

	}

	c.JSON(http.StatusOK, gin.H{

		"message": "Search restaurant success",
		"data":    restaurants,
	})
}
