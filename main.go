package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	// Create a new Gin router
	r := gin.Default()

	// Enable CORS using Gin's built-in middleware
	r.Use(cors.Default()) // This applies default CORS settings, allowing all origins

	// Define routes (example)
	r.GET("/uploadAndCompressPDF", func(c *gin.Context) {
		// Your API logic here
		c.String(200, "Hello, World!")
	})

	// Start the HTTPS server with TLS
	log.Println("Starting HTTPS server on :7070")
	err := r.RunTLS(":7070", "/etc/letsencrypt/live/34.123.178.159.nip.io/fullchain.pem", "/etc/letsencrypt/live/34.123.178.159.nip.io/privkey.pem")
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
