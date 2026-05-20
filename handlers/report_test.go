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

func TestGetBookingSummary_Success(t *testing.T) {
	// 1. สร้าง mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	database.DB = db // inject mock เข้า package

	// 2. กำหนดว่า query นี้ควร return อะไร
	rows := sqlmock.NewRows([]string{
		"total_bookings", "cancelled_bookings", "completed_bookings", "total_revenue",
	}).AddRow(10, 2, 5, 3500.00)

	mock.ExpectQuery("SELECT").
		WithArgs("1", "2026-01-01", "2026-12-31").
		WillReturnRows(rows)

	// 3. สร้าง Gin context จำลอง
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET",
		"/restaurants/1/reports/booking/summary?start_date=2026-01-01&end_date=2026-12-31",
		nil,
	)
	c.Params = gin.Params{{Key: "res_id", Value: "1"}}

	// 4. เรียก handler
	GetBookingSummary(c)

	// 5. Assert ผลลัพธ์
	assert.Equal(t, http.StatusOK, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, float64(10), result["total_bookings"])
	assert.Equal(t, float64(3500), result["total_revenue"])
}

func TestGetBookingSummary_MissingDates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/restaurants/1/reports/booking/summary", nil)
	c.Params = gin.Params{{Key: "res_id", Value: "1"}}

	GetBookingSummary(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
