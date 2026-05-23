package handlers

import (
	"bytes"
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

func TestUpdateBookingStatus_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectExec("UPDATE BOOKINGS").
		WithArgs("confirmed", "1", "1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{"status": "confirmed"})
	c.Request, _ = http.NewRequest("PATCH", "/booking/1/1/status", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "res_id", Value: "1"},
		{Key: "book_id", Value: "1"},
	}

	UpdateBookingStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateBookingStatus_InvalidStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{"status": "invalidstatus"})
	c.Request, _ = http.NewRequest("PATCH", "/booking/1/1/status", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "res_id", Value: "1"},
		{Key: "book_id", Value: "1"},
	}

	UpdateBookingStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateBookingStatus_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectExec("UPDATE BOOKINGS").
		WithArgs("confirmed", "1", "999").
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{"status": "confirmed"})
	c.Request, _ = http.NewRequest("PATCH", "/booking/1/999/status", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "res_id", Value: "1"},
		{Key: "book_id", Value: "999"},
	}

	UpdateBookingStatus(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateBookingStatus_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectExec("UPDATE BOOKINGS").
		WillReturnError(fmt.Errorf("db error"))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{"status": "confirmed"})
	c.Request, _ = http.NewRequest("PATCH", "/booking/1/1/status", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "res_id", Value: "1"},
		{Key: "book_id", Value: "1"},
	}

	UpdateBookingStatus(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateBookingStatus_MissingBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/booking/1/1/status", bytes.NewBuffer([]byte("{}")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{
		{Key: "res_id", Value: "1"},
		{Key: "book_id", Value: "1"},
	}

	UpdateBookingStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
