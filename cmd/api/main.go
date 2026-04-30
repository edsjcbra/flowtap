package main

import (
	"log"
	"os"

	"github.com/edsjcbra/flowtap/internal/database"
	"github.com/edsjcbra/flowtap/internal/handlers"
	"github.com/edsjcbra/flowtap/internal/scheduler"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	scheduler.Start()

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// CLIENTS
	router.POST("/clients", handlers.CreateClient)

	// INVOICES
	router.POST("/invoices", handlers.CreateInvoice)
	router.POST("/invoices/:id/pay", handlers.MarkAsPaid)
	router.GET("/invoices", handlers.ListInvoices)
	router.POST("/stripe/webhook", handlers.StripeWebhook)

	log.Println("API KEY:", os.Getenv("RESEND_API_KEY"))

	log.Println("Server running on :8080")
	router.Run(":8080")
}
