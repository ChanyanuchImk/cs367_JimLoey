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

func TestLogin_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	rows := sqlmock.NewRows([]string{"user_id", "password"}).
		AddRow(1, "123456")
	mock.ExpectQuery("SELECT user_id, password").
		WithArgs("jane@email.com").
		WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{
		"email":    "jane@email.com",
		"password": "123456",
	})
	c.Request, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	assert.NotEmpty(t, result["token"])
}

func TestLogin_WrongPassword(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	rows := sqlmock.NewRows([]string{"user_id", "password"}).
		AddRow(1, "123456")
	mock.ExpectQuery("SELECT user_id, password").
		WithArgs("jane@email.com").
		WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{
		"email":    "jane@email.com",
		"password": "wrongpass",
	})
	c.Request, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_UserNotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectQuery("SELECT user_id, password").
		WithArgs("notfound@email.com").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "password"}))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]string{
		"email":    "notfound@email.com",
		"password": "123456",
	})
	c.Request, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRegister_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/auth/register", bytes.NewBuffer([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	Register(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	database.DB = db

	mock.ExpectExec("INSERT INTO USERS").
		WithArgs("คุณทดสอบ", "test@email.com", "0812345678", "123456").
		WillReturnError(fmt.Errorf("duplicate entry"))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body, _ := json.Marshal(map[string]interface{}{
		"name": "คุณทดสอบ", "email": "test@email.com",
		"phone": "0812345678", "password": "123456",
	})
	c.Request, _ = http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	Register(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestLogin_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte("invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	Login(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
