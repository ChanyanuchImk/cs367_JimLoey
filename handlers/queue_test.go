package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-api/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetQueues_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	rows := sqlmock.NewRows([]string{
		"queue_id", "restaurant_id", "user_id", "queue_number", "number_of_people", "status", "created_at",
	}).
		AddRow(1, 1, 2, 1, 3, "waiting", "2026-05-20 10:00:00").
		AddRow(2, 1, 3, 2, 2, "calling", "2026-05-20 10:05:00")

	mock.ExpectQuery("SELECT queue_id, restaurant_id, user_id, queue_number, number_of_people, status, created_at FROM QUEUES").
		WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/queues", nil)

	GetQueues(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, 2, len(result))
}

func TestGetQueues_Empty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT queue_id, restaurant_id, user_id, queue_number, number_of_people, status, created_at FROM QUEUES").
		WillReturnRows(sqlmock.NewRows([]string{
			"queue_id", "restaurant_id", "user_id", "queue_number", "number_of_people", "status", "created_at",
		}))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/queues", nil)

	GetQueues(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, 0, len(result))
}

func TestGetQueues_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT").
		WillReturnError(fmt.Errorf("db error"))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/queues", nil)

	GetQueues(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
