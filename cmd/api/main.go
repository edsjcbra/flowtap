package main

import (
	"log"
	"os"

	"github.com/edsjcbra/flowtap/internal/database"
	"github.com/edsjcbra/flowtap/internal/handlers"
	"github.com/edsjcbra/flowtap/internal/middleware"
	"github.com/edsjcbra/flowtap/internal/worker"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	database.Connect()
	log.Println("RESEND KEY:", os.Getenv("RESEND_API_KEY"))
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
	auth.GET("/clients", handlers.ListClients)
	auth.POST("/invoices", handlers.CreateInvoice)
	auth.GET("/invoices", handlers.ListInvoices)
	auth.POST("/invoices/:id/pay", handlers.MarkAsPaid)
	auth.POST("/invoices/:id/cancel", handlers.CancelInvoice)
	auth.PUT("/clients/:id", handlers.UpdateClient)
	auth.DELETE("/clients/:id", handlers.DeleteClient)

	auth.DELETE("/invoices/:id", handlers.DeleteInvoice)

	log.Println("API KEY:", os.Getenv("RESEND_API_KEY"))

	log.Println("Server running on :8080")

	go worker.StartWorker()
	router.Run(":8080")

}
