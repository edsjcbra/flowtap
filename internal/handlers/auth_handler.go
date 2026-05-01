package handlers

import (
	"github.com/edsjcbra/flowtap/internal/database"
	"github.com/edsjcbra/flowtap/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(c *gin.Context) {
	var req AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	_, err := database.DB.Exec(`
		INSERT INTO users (email, password)
		VALUES ($1, $2)
	`, req.Email, req.Password)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "user created"})
}

func Login(c *gin.Context) {
	var req AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	var userID int
	var storedPassword string

	err := database.DB.QueryRow(`
		SELECT id, password FROM users WHERE email = $1
	`, req.Email).Scan(&userID, &storedPassword)

	if err != nil {
		c.JSON(401, gin.H{"error": "user not found"})
		return
	}

	if storedPassword != req.Password {
		c.JSON(401, gin.H{"error": "wrong password"})
		return
	}

	token, err := services.GenerateToken(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
