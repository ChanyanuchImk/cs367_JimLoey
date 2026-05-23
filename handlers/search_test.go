package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-api/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSearchRestaurants_NoKeyword(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	rows := sqlmock.NewRows([]string{
		"restaurant_id", "name", "location", "phone", "open_time", "close_time",
	}).
		AddRow(1, "ร้านข้าวแกง", "รังสิต", "0811111111", "08:00:00", "16:00:00").
		AddRow(2, "Sushi House", "ฟิวเจอร์พาร์ค", "0822222222", "10:00:00", "22:00:00")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/restaurants", nil)

	SearchRestaurants(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "Search restaurant success", result["message"])
}

func TestSearchRestaurants_WithKeyword(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	rows := sqlmock.NewRows([]string{
		"restaurant_id", "name", "location", "phone", "open_time", "close_time",
	}).AddRow(2, "Sushi House", "ฟิวเจอร์พาร์ค", "0822222222", "10:00:00", "22:00:00")

	mock.ExpectQuery("SELECT").WithArgs("%Sushi%").WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/restaurants?keyword=Sushi", nil)

	SearchRestaurants(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetRestaurantDetail_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	rows := sqlmock.NewRows([]string{
		"restaurant_id", "name", "description", "location", "phone",
		"open_time", "close_time", "avg_rating", "total_reviews",
	}).AddRow(1, "ร้านข้าวแกง", "อร่อยมาก", "รังสิต", "0811111111", "08:00:00", "16:00:00", 4.5, 10)

	mock.ExpectQuery("SELECT").WithArgs("1").WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/restaurants/1", nil)
	c.Params = gin.Params{{Key: "res_id", Value: "1"}}

	GetRestaurantDetail(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, "Get restaurant detail success", result["message"])
}

func TestGetRestaurantDetail_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT").WithArgs("999").
		WillReturnRows(sqlmock.NewRows([]string{
			"restaurant_id", "name", "description", "location", "phone",
			"open_time", "close_time", "avg_rating", "total_reviews",
		}))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/restaurants/999", nil)
	c.Params = gin.Params{{Key: "res_id", Value: "999"}}

	GetRestaurantDetail(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
