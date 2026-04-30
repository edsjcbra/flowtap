package main

import (
	"log"

	"github.com/edsjcbra/flowtap/internal/database"
	"github.com/edsjcbra/flowtap/internal/handlers"
	"github.com/edsjcbra/flowtap/internal/scheduler"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	scheduler.Start()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.POST("/clients", handlers.CreateClient)
	router.POST("/invoices", handlers.CreateInvoice)

	log.Println("Server running on :8080")
	router.Run(":8080")
}
