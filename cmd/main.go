package main

import (
	"context"
	"fmt"
	"gcpwebserv"
	"log/slog"
	"os"
)

func main() {
	// Create a context
	ctx := context.Background()

	//Configures a default http server with the GCP default such as port 8080.
	// It can be overridden by providing a preconfigured http.Server
	server := gcpwebserv.Setup(nil)
	err := gcpwebserv.SetupStaticFileHandler("/static/", "static")
	if err != nil {
		os.Exit(1)
	}

	err = gcpwebserv.RegisterGCSHandler(ctx)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	err = gcpwebserv.RegisterFirestoreNativeHandler(ctx, "ups-hacka", "(default)")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	err = gcpwebserv.RegisterPubSubHandler(ctx, "ups-hacka", "poop")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	if err := gcpwebserv.Run(ctx, server); err != nil {
		slog.Info("Server Exiting", "status", err)
	}

}
