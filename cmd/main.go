package main

import (
	"context"
	"log/slog"
	"os"

	gcpwebserver "github.com/ct-ldeakm/gcpwebserver"
)

func main() {
	// Create a context
	ctx := context.Background()

	//Configures a default http server with the GCP default such as port 8080.
	// It can be overridden by providing a pre configured http.Server
	server := gcpwebserver.Setup(nil)
	err := gcpwebserver.SetupStaticFileHandler("/static/", "static")
	if err != nil {
		os.Exit(1)
	}

	if err := gcpwebserver.Run(ctx, server); err != nil {
		slog.Info("Server Exiting", "status", err)
	}

}
