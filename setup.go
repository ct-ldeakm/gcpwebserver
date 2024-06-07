package gcpwebserv

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func Setup(server *http.Server, prebuilt ...http.HandlerFunc) *http.Server {

	// Setting some default timeouts based on the intent of hosting GRCP or rest services
	// and interacting with GCP services like BQ and Spanner from Cloud Run.
	if server == nil {
		server = &http.Server{
			IdleTimeout:       time.Minute * 1,
			ReadTimeout:       time.Minute * 15,
			ReadHeaderTimeout: time.Second * 30,
			WriteTimeout:      time.Minute * 60,
		}

		// The port can be set using ENV PORT but will be overridden if
		// the service object is provided to setup.
		port, isSet := os.LookupEnv("PORT")

		if isSet {
			server.Addr = fmt.Sprintf(":%s", port)
		}
	}

	// Setting the default port to 8080 incase nothing was set.
	if server.Addr == "" {
		server.Addr = ":8080"
	}

	return server

}
