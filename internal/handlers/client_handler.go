// Package handlers contain the handlers of the application domains
package handlers

import (
	"net/http"

	"github.com/edsjcbra/flowtap/internal/services"
	"github.com/gin-gonic/gin"
)

type CreateClientRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func CreateClient(c *gin.Context) {
	var req CreateClientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.CreateClient(req.Name, req.Email, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "client created"})
}
