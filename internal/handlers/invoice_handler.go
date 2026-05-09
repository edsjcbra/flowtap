package handlers

import (
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

// ================= CREATE =================

func CreateInvoice(c *gin.Context) {

	var req CreateInvoiceRequest

	userID := c.GetInt("user_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	dueDate, err := time.Parse(
		time.RFC3339,
		req.DueDate,
	)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid date",
		})
		return
	}

	invoiceID, err := services.CreateInvoice(
		req.ClientID,
		req.Amount,
		dueDate,
		userID,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// stripe link

	url, err := services.CreateCheckoutSession(
		req.Amount,
		invoiceID,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// save payment url

	_, err = database.DB.Exec(`
		UPDATE invoices
		SET payment_url = $1
		WHERE id = $2
	`,
		url,
		invoiceID,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"invoice_id":  invoiceID,
		"payment_url": url,
	})
}

// ================= LIST =================

func ListInvoices(c *gin.Context) {

	userID := c.GetInt("user_id")

	rows, err := database.DB.Query(`
		SELECT
			i.id,
			i.amount,
			i.status,
			i.due_date,
			c.name
		FROM invoices i
		JOIN clients c
			ON i.client_id = c.id
		WHERE i.user_id = $1
		ORDER BY i.id DESC
	`, userID)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer rows.Close()

	var invoices []gin.H

	for rows.Next() {

		var id int
		var amount float64
		var status string
		var dueDate time.Time
		var clientName string

		err := rows.Scan(
			&id,
			&amount,
			&status,
			&dueDate,
			&clientName,
		)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		invoices = append(invoices, gin.H{
			"id":          id,
			"client_name": clientName,
			"amount":      amount,
			"status":      status,
			"due_date":    dueDate,
		})
	}

	c.JSON(200, invoices)
}

// ================= MARK AS PAID =================

func MarkAsPaid(c *gin.Context) {

	idStr := c.Param("id")

	invoiceID, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid id",
		})
		return
	}

	err = services.MarkInvoiceAsPaid(invoiceID)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "paid",
	})
}

// ================= CANCEL =================

func CancelInvoice(c *gin.Context) {

	idStr := c.Param("id")

	invoiceID, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid id",
		})
		return
	}

	err = services.CancelInvoice(invoiceID)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "cancelled",
	})
}

// ================= DELETE =================

func DeleteInvoice(c *gin.Context) {

	idStr := c.Param("id")

	invoiceID, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid id",
		})
		return
	}

	err = services.DeleteInvoice(invoiceID)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "deleted",
	})
}
