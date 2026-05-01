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
	userID := c.GetInt("user_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	id, err := services.CreateInvoice(req.ClientID, req.Amount, dueDate, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invoice_id": id})
}

func MarkAsPaid(c *gin.Context) {
	idStr := c.Param("id")
	userID := c.GetInt("user_id")

	invoiceID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invoice id"})
		return
	}


	var exists bool
	err = database.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM invoices
			WHERE id = $1 AND user_id = $2
		)
	`, invoiceID, userID).Scan(&exists)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.JSON(403, gin.H{"error": "not allowed"})
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
	userID := c.GetInt("user_id")

	rows, err := database.DB.Query(`
		SELECT id, client_id, amount, status
		FROM invoices
		WHERE user_id = $1
		ORDER BY id DESC
	`, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	invoices := []gin.H{}

	for rows.Next() {
		var id, clientID int
		var amount float64
		var status string

		if err := rows.Scan(&id, &clientID, &amount, &status); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		invoices = append(invoices, gin.H{
			"id":        id,
			"client_id": clientID,
			"amount":    amount,
			"status":    status,
		})
	}

	c.JSON(200, invoices)
}