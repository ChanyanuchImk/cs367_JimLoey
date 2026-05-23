package handlers

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"
	"restaurant-api/database"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTableStatus_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectExec("UPDATE TABLES").
		WithArgs("occupied", "1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{"status": "occupied"})
	c.Request, _ = http.NewRequest("PATCH", "/tables/1/status", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "table_id", Value: "1"}}

	UpdateTableStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateTableStatus_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectExec("UPDATE TABLES").
		WithArgs("occupied", "999").
		WillReturnResult(sqlmock.NewResult(0, 0))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{"status": "occupied"})
	c.Request, _ = http.NewRequest("PATCH", "/tables/999/status", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "table_id", Value: "999"}}

	UpdateTableStatus(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetTableCount_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT COUNT").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/tables/count", nil)

	GetTableCount(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.Equal(t, float64(5), result["total_tables"])
}

func TestUpdateTableStatus_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PATCH", "/tables/1/status", bytes.NewBuffer([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "table_id", Value: "1"}}

	UpdateTableStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTableCount_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT COUNT").
		WillReturnError(fmt.Errorf("db error"))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/tables/count", nil)

	GetTableCount(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
