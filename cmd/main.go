package main

import (
	"context"
	"gcpwebserv"
	"log"
)

func main() {
	// Create a context
	ctx := context.Background()

	//Configures a default http server with the GCP default such as port 8080.
	// It can be overridden by providing a preconfigured http.Server
	server := gcpwebserv.Setup(nil)
	gcpwebserv.SetupStaticFileHandler("/static/", "static/")
	// Add Static Routes
	//gcpwebserv.SetupStaticFileHandler("/app/", "./app")
	//addRoutes()

	if err := gcpwebserv.Run(ctx, server); err != nil {
		log.Printf("%s", err)
	}

}
