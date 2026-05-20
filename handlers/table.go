package handlers

import (
	"net/http"
	"restaurant-api/database"

	"github.com/gin-gonic/gin"
)

type UpdateTableStatusRequest struct {
	Status string `json:"status"`
}

func UpdateTableStatus(c *gin.Context) {

	tableID := c.Param("table_id")

	var req UpdateTableStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})

		return
	}

	query := `
	UPDATE TABLES
	SET status = ?
	WHERE table_id = ?
	`

	result, err := database.DB.Exec(
		query,
		req.Status,
		tableID,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "table not found",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "table status updated",
	})
}

func GetTableCount(c *gin.Context) {

	var count int

	query := `
	SELECT COUNT(*)
	FROM TABLES
	`

	err := database.DB.QueryRow(query).Scan(&count)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_tables": count,
	})
}