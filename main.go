package main

import (
	"testproject/routes"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {

	r := routes.InitializeRoutes()
	log.Println("Starting HTTPS server on :443")
	err := http.ListenAndServeTLS(
		":443", 
		"/etc/letsencrypt/live/34.123.178.159.nip.io/fullchain.pem", 
		"/etc/letsencrypt/live/34.123.178.159.nip.io/privkey.pem", 
		nil,
	)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
