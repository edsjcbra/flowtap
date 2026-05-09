package handlers

import (
	"net/http"
	"strconv"

	"github.com/edsjcbra/flowtap/internal/database"
	"github.com/edsjcbra/flowtap/internal/services"
	"github.com/gin-gonic/gin"
)

// ================= CREATE CLIENT =================

type CreateClientRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateClient(c *gin.Context) {
	var req CreateClientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 🔥 pega user_id do token
	userID := c.GetInt("user_id")

	_, err := database.DB.Exec(`
		INSERT INTO clients (name, email, user_id)
		VALUES ($1, $2, $3)
	`, req.Name, req.Email, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "client created",
	})
}

// ================= LIST CLIENTS =================

func ListClients(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := database.DB.Query(`
		SELECT id, name, email
		FROM clients
		WHERE user_id = $1
		ORDER BY id DESC
	`, userID)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var clients []gin.H

	for rows.Next() {
		var id int
		var name, email string

		rows.Scan(&id, &name, &email)

		clients = append(clients, gin.H{
			"id":    id,
			"name":  name,
			"email": email,
		})
	}

	c.JSON(200, clients)
}

func UpdateClient(c *gin.Context) {
	idStr := c.Param("id")
	userID := c.GetInt("user_id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	err = services.UpdateClient(id, body.Name, body.Email, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "updated"})
}

func DeleteClient(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := services.DeleteClient(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "deleted"})
}


