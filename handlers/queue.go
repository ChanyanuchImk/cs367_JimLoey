package handlers

import (
	"net/http"
	"restaurant-api/database"

	"github.com/gin-gonic/gin"
)

func GetQueues(c *gin.Context) {

	rows, err := database.DB.Query(`
		SELECT
			queue_id,
			queue_number,
			status
		FROM QUEUES
	`)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	defer rows.Close()

	type Queue struct {
		QueueID     int    `json:"queue_id"`
		QueueNumber int    `json:"queue_number"`
		Status      string `json:"status"`
	}

	var queues []Queue

	for rows.Next() {

		var q Queue

		rows.Scan(
			&q.QueueID,
			&q.QueueNumber,
			&q.Status,
		)

		queues = append(queues, q)
	}

	c.JSON(http.StatusOK, queues)
}
