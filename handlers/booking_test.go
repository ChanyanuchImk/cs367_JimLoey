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

// ========== CreateBooking ==========

func TestCreateBooking_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	// mock เช็ค restaurant
	mock.ExpectQuery("SELECT restaurant_id FROM RESTAURANTS").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"restaurant_id"}).AddRow(1))

	// mock หา table ว่าง
	mock.ExpectQuery("SELECT table_id FROM TABLES").
		WithArgs(1, 2, 1, "2026-06-01", "12:00:00").
		WillReturnRows(sqlmock.NewRows([]string{"table_id"}).AddRow(1))

	// mock insert booking
	mock.ExpectExec("INSERT INTO BOOKINGS").
		WithArgs(1, 1, 1, "2026-06-01", "12:00:00", "14:00:00", 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]interface{}{
		"user_id":          1,
		"restaurant_id":    1,
		"booking_date":     "2026-06-01",
		"start_time":       "12:00:00",
		"number_of_people": 2,
	})
	c.Request, _ = http.NewRequest("POST", "/booking", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateBooking(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "pending", result["status"])
}

func TestCreateBooking_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]interface{}{
		"user_id": 1,
		// ขาด restaurant_id, booking_date, start_time, number_of_people
	})
	c.Request, _ = http.NewRequest("POST", "/booking", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateBooking(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateBooking_NoTableAvailable(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT restaurant_id FROM RESTAURANTS").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"restaurant_id"}).AddRow(1))

	// ไม่มี table ว่าง
	mock.ExpectQuery("SELECT table_id FROM TABLES").
		WithArgs(1, 2, 1, "2026-06-01", "12:00:00").
		WillReturnRows(sqlmock.NewRows([]string{"table_id"}))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]interface{}{
		"user_id":          1,
		"restaurant_id":    1,
		"booking_date":     "2026-06-01",
		"start_time":       "12:00:00",
		"number_of_people": 2,
	})
	c.Request, _ = http.NewRequest("POST", "/booking", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateBooking(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

// ========== GetBookings ==========

func TestGetBookings_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	rows := sqlmock.NewRows([]string{
		"booking_id", "user_id", "restaurant_id", "table_id",
		"booking_date", "start_time", "end_time",
		"number_of_people", "total_price", "status", "created_at",
	}).AddRow(1, 1, 1, 1, "2026-06-01", "12:00:00", "14:00:00", 2, 500.00, "confirmed", "2026-05-01 10:00:00")

	mock.ExpectQuery("SELECT").WithArgs("1").WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/booking/1", nil)
	c.Params = gin.Params{{Key: "res_id", Value: "1"}}

	GetBookings(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, 1, len(result))
}

func TestGetBookings_Empty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT").WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{
			"booking_id", "user_id", "restaurant_id", "table_id",
			"booking_date", "start_time", "end_time",
			"number_of_people", "total_price", "status", "created_at",
		}))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/booking/1", nil)
	c.Params = gin.Params{{Key: "res_id", Value: "1"}}

	GetBookings(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, 0, len(result))
}

func TestGetBookingsByUser_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	rows := sqlmock.NewRows([]string{
		"booking_id", "user_id", "restaurant_id", "table_id",
		"booking_date", "start_time", "end_time",
		"number_of_people", "total_price", "status", "created_at",
	}).AddRow(1, 1, 1, 1, "2026-06-01", "12:00:00", "14:00:00", 2, 500.00, "confirmed", "2026-05-01 10:00:00")

	mock.ExpectQuery("SELECT").WithArgs("1").WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/booking/user/1", nil)
	c.Params = gin.Params{{Key: "user_id", Value: "1"}}

	GetBookingsByUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, 1, len(result))
}

func TestGetBookingsByRestaurant_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	rows := sqlmock.NewRows([]string{
		"booking_id", "user_id", "restaurant_id", "table_id",
		"booking_date", "start_time", "end_time",
		"number_of_people", "total_price", "status", "created_at",
	}).AddRow(1, 1, 1, 1, "2026-06-01", "12:00:00", "14:00:00", 2, 500.00, "confirmed", "2026-05-01 10:00:00")

	mock.ExpectQuery("SELECT").WithArgs("1").WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/booking/restaurant/1", nil)
	c.Params = gin.Params{{Key: "res_id", Value: "1"}}

	GetBookingsByRestaurant(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, 1, len(result))
}

func TestCreateBooking_InvalidTime(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]interface{}{
		"user_id": 1, "restaurant_id": 1,
		"booking_date": "2026-06-01", "start_time": "invalid-time",
		"number_of_people": 2,
	})
	c.Request, _ = http.NewRequest("POST", "/booking", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateBooking(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateBooking_RestaurantNotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT restaurant_id FROM RESTAURANTS").
		WithArgs(999).
		WillReturnRows(sqlmock.NewRows([]string{"restaurant_id"}))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]interface{}{
		"user_id": 1, "restaurant_id": 999,
		"booking_date": "2026-06-01", "start_time": "12:00:00",
		"number_of_people": 2,
	})
	c.Request, _ = http.NewRequest("POST", "/booking", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateBooking(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetBookingsByUser_Empty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT").WithArgs("99").
		WillReturnRows(sqlmock.NewRows([]string{
			"booking_id", "user_id", "restaurant_id", "table_id",
			"booking_date", "start_time", "end_time",
			"number_of_people", "total_price", "status", "created_at",
		}))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/booking/user/99", nil)
	c.Params = gin.Params{{Key: "user_id", Value: "99"}}

	GetBookingsByUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetBookingsByRestaurant_Empty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT").WithArgs("99").
		WillReturnRows(sqlmock.NewRows([]string{
			"booking_id", "user_id", "restaurant_id", "table_id",
			"booking_date", "start_time", "end_time",
			"number_of_people", "total_price", "status", "created_at",
		}))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/booking/restaurant/99", nil)
	c.Params = gin.Params{{Key: "res_id", Value: "99"}}

	GetBookingsByRestaurant(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetBookings_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT").WithArgs("1").
		WillReturnError(fmt.Errorf("db error"))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/booking/1", nil)
	c.Params = gin.Params{{Key: "res_id", Value: "1"}}

	GetBookings(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetBookingsByUser_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT").WithArgs("1").
		WillReturnError(fmt.Errorf("db error"))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/booking/user/1", nil)
	c.Params = gin.Params{{Key: "user_id", Value: "1"}}

	GetBookingsByUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSearchRestaurants_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT").
		WillReturnError(fmt.Errorf("db error"))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/restaurants", nil)

	SearchRestaurants(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetRestaurantDetail_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT").WithArgs("1").
		WillReturnError(fmt.Errorf("db error"))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/restaurants/1", nil)
	c.Params = gin.Params{{Key: "res_id", Value: "1"}}

	GetRestaurantDetail(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateBooking_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/booking", bytes.NewBuffer([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	CreateBooking(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetBookingsByRestaurant_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT").WithArgs("1").
		WillReturnError(fmt.Errorf("db error"))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/booking/restaurant/1", nil)
	c.Params = gin.Params{{Key: "res_id", Value: "1"}}

	GetBookingsByRestaurant(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
