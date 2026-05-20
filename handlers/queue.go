package handlers

import (
	"net/http"
	"restaurant-api/database"

	"github.com/gin-gonic/gin"
)

type Queue struct {
	QueueID        int    `json:"queue_id"`
	RestaurantID   int    `json:"restaurant_id"`
	UserID         int    `json:"user_id"`
	QueueNumber    int    `json:"queue_number"`
	NumberOfPeople int    `json:"number_of_people"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
}

func GetQueues(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT queue_id, restaurant_id, user_id, queue_number, number_of_people, status, created_at
		FROM QUEUES
		ORDER BY created_at ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query queues"})
		return
	}
	defer rows.Close()

	var queues []Queue
	for rows.Next() {
		var q Queue
		err := rows.Scan(
			&q.QueueID,
			&q.RestaurantID,
			&q.UserID,
			&q.QueueNumber,
			&q.NumberOfPeople,
			&q.Status,
			&q.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse queue data"})
			return
		}
		queues = append(queues, q)
	}

	if queues == nil {
		queues = []Queue{}
	}

	c.JSON(http.StatusOK, queues)
}
