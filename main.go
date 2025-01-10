package main

import (
	"log"
	"github.com/gin-contrib/cors"
	"testproject/routes" // Import your routes package
)

func main() {
	// Initialize routes from the routes package
	r := routes.InitializeRoutes()

	// Enable CORS using Gin's built-in middleware
    	// Enable CORS using Gin's built-in middleware
   	 r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
	    	"*", // Allow all origins (use for testing)
	    	"https://34.123.178.159.nip.io",  // Your API domain
		},
    	})) // This applies default CORS settings, allowing all origins

	// Start the HTTPS server with TLS
	log.Println("Starting HTTPS server on :7070")
	err := r.RunTLS(":7070", "/etc/letsencrypt/live/34.123.178.159.nip.io/fullchain.pem", "/etc/letsencrypt/live/34.123.178.159.nip.io/privkey.pem")
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
