// Package handlers deals with the http requests
package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/edsjcbra/flowtap/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
)

func StripeWebhook(c *gin.Context) {
	log.Println("🔥 WEBHOOK RECEIVED")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "read error"})
		return
	}

	var event stripe.Event

	err = json.Unmarshal(body, &event)
	if err != nil {
		log.Println("JSON error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	log.Println("Event type:", event.Type)

	if event.Type == "checkout.session.completed" {

		var session stripe.CheckoutSession

		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Println("JSON error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		log.Println("Metadata:", session.Metadata)

		invoiceIDStr := session.Metadata["invoice_id"]
		invoiceID, _ := strconv.Atoi(invoiceIDStr)

		log.Println("Marking invoice as paid:", invoiceID)

		err = services.MarkInvoiceAsPaid(invoiceID)
		if err != nil {
			log.Println("DB error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
