package gcpwebserv

import (
	"fmt"
	"log/slog"
	"net/http"
)

func RegisterDefaultHandler() error {
	Route("GET /", defaultDump)
	return nil
}

func defaultDump(w http.ResponseWriter, r *http.Request) {
	slog.Info("In Default")

	for k, v := range r.Header {
		w.Write([]byte(fmt.Sprintf("%s:%s\n", k, v)))
	}
	w.Write([]byte(fmt.Sprintf("%s", r.URL.Path)))

}
