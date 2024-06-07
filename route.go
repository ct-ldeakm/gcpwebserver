package gcpwebserv

import "net/http"

// Adding a global mux here to simplify
var mux *http.ServeMux

func init() {
	mux = http.NewServeMux()
}

// Add your handlers here using the mux

// Route takes pefined handler function attaches the path route and adds it
// to the serving multiplixer
func Route(path string, handler http.HandlerFunc) {
	mux.HandleFunc(path, handler)
}

func RouteHandler(path string, handler http.Handler) {
	mux.Handle(path, handler)
}
