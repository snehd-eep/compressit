package main

import (
	"log"
	"net/http"
	"github.com/rs/cors"
	"testproject/routes"
)

func main() {
	// Initialize routes
	r := routes.InitializeRoutes()

	// Enable CORS for the routes
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins (adjust this for production)
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed methods
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"}, // Allowed headers
	})

	// Wrap the router with CORS middleware
	r = corsHandler.Handler(r)

	// Start the HTTPS server
	log.Println("Starting HTTPS server on :7070")
	err := http.ListenAndServeTLS(
		":7070", // Correct port here
		"/etc/letsencrypt/live/34.123.178.159.nip.io/fullchain.pem", // Path to the certificate
		"/etc/letsencrypt/live/34.123.178.159.nip.io/privkey.pem",  // Path to the private key
		r, // Pass the router as the handler
	)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
