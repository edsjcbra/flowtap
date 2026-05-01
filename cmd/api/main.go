package main

import (
	"log"
	"os"

	"github.com/edsjcbra/flowtap/internal/database"
	"github.com/edsjcbra/flowtap/internal/handlers"
	"github.com/edsjcbra/flowtap/internal/middleware"
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

	// 🔓 rotas públicas
	router.POST("/signup", handlers.Signup)
	router.POST("/login", handlers.Login)
	router.POST("/stripe/webhook", handlers.StripeWebhook)

	// 🔒 rotas protegidas
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.POST("/clients", handlers.CreateClient)
	auth.POST("/invoices", handlers.CreateInvoice)
	auth.GET("/invoices", handlers.ListInvoices)
	auth.POST("/invoices/:id/pay", handlers.MarkAsPaid)

	log.Println("API KEY:", os.Getenv("RESEND_API_KEY"))

	log.Println("Server running on :8080")
	router.Run(":8080")
}
