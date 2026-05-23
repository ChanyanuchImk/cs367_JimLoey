package handlers

import (
	"database/sql"
	"errors"
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
		FROM RESTAURANTS
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

func GetRestaurantDetail(c *gin.Context) {

	id := c.Param("res_id")

	query := `
	SELECT
		r.restaurant_id,
		r.name,
		COALESCE(r.description, '') AS description,
		r.location,
		r.phone,
		r.open_time,
		r.close_time,
		COALESCE(ROUND(AVG(rv.rating),1), 0) AS avg_rating,
		COUNT(rv.review_id) AS total_reviews
	FROM RESTAURANTS r
	LEFT JOIN REVIEWS rv
		ON r.restaurant_id = rv.restaurant_id
	WHERE r.restaurant_id = ?
	GROUP BY r.restaurant_id
`

	type RestaurantDetail struct {
		RestaurantID int     `json:"restaurant_id"`
		Name         string  `json:"name"`
		Description  string  `json:"description"`
		Location     string  `json:"location"`
		Phone        string  `json:"phone"`
		OpenTime     string  `json:"open_time"`
		CloseTime    string  `json:"close_time"`
		AvgRating    float64 `json:"avg_rating"`
		TotalReviews int     `json:"total_reviews"`
	}

	var restaurant RestaurantDetail

	err := database.DB.QueryRow(query, id).Scan(
		&restaurant.RestaurantID,
		&restaurant.Name,
		&restaurant.Description,
		&restaurant.Location,
		&restaurant.Phone,
		&restaurant.OpenTime,
		&restaurant.CloseTime,
		&restaurant.AvgRating,
		&restaurant.TotalReviews,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Restaurant not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get restaurant detail success",
		"data":    restaurant,
	})
}
