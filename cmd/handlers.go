package main

import (
	"gcpwebserv"
	"net/http"
)

// addRoutes simpley wraps the lig function to add a route to the mux handler
// so that a one to one mapping of route to handler can be seen for troubleshooting
func addRoutes() {
	gcpwebserv.Route("/{id}", pathTest)
}

// Create you handelers below this line
func pathTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.PathValue("id")))
}
