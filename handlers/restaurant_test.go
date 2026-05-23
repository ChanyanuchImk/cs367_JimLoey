package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"restaurant-api/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateRestaurant_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectExec("INSERT INTO RESTAURANTS").
		WithArgs(4, "ร้านทดสอบ", "อร่อยมาก", "กรุงเทพฯ", "0811111111", "08:00:00", "20:00:00").
		WillReturnResult(sqlmock.NewResult(1, 1))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]interface{}{
		"owner_id":    4,
		"name":        "ร้านทดสอบ",
		"description": "อร่อยมาก",
		"location":    "กรุงเทพฯ",
		"phone":       "0811111111",
		"open_time":   "08:00:00",
		"close_time":  "20:00:00",
	})
	c.Request, _ = http.NewRequest("POST", "/restaurants", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateRestaurant(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateRestaurant_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/restaurants", bytes.NewBuffer([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateRestaurant(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectExec("INSERT INTO USERS").
		WithArgs("คุณทดสอบ", "test@email.com", "0812345678", "123456").
		WillReturnResult(sqlmock.NewResult(1, 1))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]interface{}{
		"name":     "คุณทดสอบ",
		"email":    "test@email.com",
		"phone":    "0812345678",
		"password": "123456",
	})
	c.Request, _ = http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	Register(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateRestaurant_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectExec("INSERT INTO RESTAURANTS").
		WillReturnError(fmt.Errorf("db error"))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]interface{}{
		"owner_id": 4, "name": "ร้านทดสอบ",
		"description": "อร่อยมาก", "location": "กรุงเทพฯ",
		"phone": "0811111111", "open_time": "08:00:00", "close_time": "20:00:00",
	})
	c.Request, _ = http.NewRequest("POST", "/restaurants", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateRestaurant(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
