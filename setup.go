package gcpwebserv

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Setup configures or pass thru a preconfigured pointer to an http server.
// The defaults provided are standard set based on best use for Cloud Run.
// Port 8080 is also set as the default but can be overwritten by the PORT
// environment variable.
func Setup(server *http.Server, prebuilt ...http.HandlerFunc) *http.Server {

	// Setting some default timeouts based on the intent of hosting GRCP or rest services
	// and interacting with GCP services like BQ and Spanner from Cloud Run.
	if server == nil {
		slog.Info("Setting up default http server")
		server = &http.Server{
			IdleTimeout:       time.Minute * 1,
			ReadTimeout:       time.Minute * 15,
			ReadHeaderTimeout: time.Second * 30,
			WriteTimeout:      time.Minute * 60,
		}

	}

	// The port can be set using ENV PORT but will be overridden if
	// the service object is provided to setup.
	port, isSet := os.LookupEnv("PORT")

	if isSet {
		slog.Info("Setting provided port", "port", port)
		server.Addr = fmt.Sprintf(":%s", port)
	} else {
		// Setting the default port to 8080 incase nothing was set.
		if server.Addr == "" {
			slog.Info("Setting default port ", "port", ":8080")
			server.Addr = ":8080"
		}
	}

	return server

}
