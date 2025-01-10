package main

import (
	"log"
	"net/http"
	"testproject/routes"
	"github.com/rs/cors"
)

func main() {
	// Initialize routes
	r := routes.InitializeRoutes()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins (adjust this for production)
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allowed methods
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With"}, // Allowed headers
	})

	r = corsHandler.Handler(r)

	// Start the HTTPS server
	log.Println("Starting HTTPS server on :7070")
	err := http.ListenAndServeTLS(
		":7070",
		"/etc/letsencrypt/live/34.123.178.159.nip.io/fullchain.pem",
		"/etc/letsencrypt/live/34.123.178.159.nip.io/privkey.pem",
		r, // Pass the router as the handler
	)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
