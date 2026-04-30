package handlers

import (
	"net/http"
	"time"

	"github.com/edsjcbra/flowtap/internal/services"

	"github.com/gin-gonic/gin"
)

type CreateInvoiceRequest struct {
	ClientID int     `json:"client_id"`
	Amount   float64 `json:"amount"`
	DueDate  string  `json:"due_date"`
}

func CreateInvoice(c *gin.Context) {
	var req CreateInvoiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dueDate, _ := time.Parse(time.RFC3339, req.DueDate)

	id, err := services.CreateInvoice(req.ClientID, req.Amount, dueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invoice_id": id})
}