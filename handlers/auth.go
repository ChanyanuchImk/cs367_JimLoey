package handlers

import (
	"net/http"
	"os"
	"restaurant-api/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {

	var req LoginRequest

	err := c.ShouldBindJSON(&req)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})

		return
	}

	var userID int
	var password string

	query := `
	SELECT user_id, password
	FROM USERS
	WHERE email = ?
	`

	err = database.DB.QueryRow(
		query,
		req.Email,
	).Scan(&userID, &password)

	if err != nil {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})

		return
	}

	if req.Password != password {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Wrong password",
		})

		return
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp": time.Now().Add(
				time.Hour * 24,
			).Unix(),
		},
	)

	tokenString, _ := token.SignedString(jwtKey)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
