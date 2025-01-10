package main

import (
	"log"
	"net/http"
	"testproject/routes"
)

func main() {
	// Initialize routes
	r := routes.InitializeRoutes()

	// Start the HTTPS server
	log.Println("Starting HTTPS server on :443")
	err := http.ListenAndServeTLS(
		":443",
		"/etc/letsencrypt/live/34.123.178.159.nip.io/fullchain.pem",
		"/etc/letsencrypt/live/34.123.178.159.nip.io/privkey.pem",
		r, // Pass the router as the handler
	)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
