package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/edsjcbra/flowtap/internal/database"
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

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	id, err := services.CreateInvoice(req.ClientID, req.Amount, dueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invoice_id": id})
}

func MarkAsPaid(c *gin.Context) {
	idStr := c.Param("id")

	invoiceID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoice id"})
		return
	}

	err = services.MarkInvoiceAsPaid(invoiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "invoice marked as paid",
	})
}

func ListInvoices(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, client_id, amount, status
		FROM invoices
		ORDER BY id DESC
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var invoices []gin.H

	for rows.Next() {
		var id, clientID int
		var amount float64
		var status string

		rows.Scan(&id, &clientID, &amount, &status)

		invoices = append(invoices, gin.H{
			"id":        id,
			"client_id": clientID,
			"amount":    amount,
			"status":    status,
		})
	}

	c.JSON(200, invoices)
}
