package gcpwebserv

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var firestoreClient *firestore.Client

// RegisterFirestoreNativeHandler setup the firestore handler and caches
// the client. If the databaseId is set to "" the (default)
// databaseId for the project is used (default). This this handler will
// not work for firestore in datastore mode. The URL for this route
// is in the form https://domain.com/firestore/projectid/databaseId/collection/docid
func RegisterFirestoreNativeHandler(ctx context.Context, projectId string, databaseId string, opts ...option.ClientOption) error {
	// Set the default if the databaseId is set to ""
	if databaseId == "" {
		databaseId = "(default)"
	}

	var err error
	firestoreClient, err = firestore.NewClientWithDatabase(ctx, projectId, databaseId, opts...)
	if err != nil {
		slog.Error("Error Creating firestore client", "error", err)
		return err
	}

	clientManager.Add(fmt.Sprintf("firestore-%s-%s", projectId, databaseId), firestoreClient)

	Route("/firestore/", getDocFromFirestore)
	return nil

}

func getDocFromFirestore(w http.ResponseWriter, r *http.Request) {
	// Clean the path of the url
	path := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	log.Printf("%v", path)
	// Expected pattern is firestore/projectid/collection/docid

	// Split the path and verify min length to make the request
	sPath := strings.Split(path, "/")
	log.Printf("%v", sPath)
	if len(sPath) < 4 {
		http.Error(w, "Bad Path", http.StatusBadRequest)
		return
	}

	// Rejoin the path getting just the file path
	obj := strings.Join(sPath[3:len(sPath)], "/")
	log.Printf("%v", obj)

	// Get the forestore client from cache. This is done so we can
	// register multiple connections to mulitple projects and databased
	// with the same handler and still leverage client caching
	firestoreClient, ok := clientManager.Get(fmt.Sprintf("firestore-%s-%s", sPath[1], sPath[2])).(*firestore.Client)
	if !ok {
		slog.Error("Error getting firestore client from cache", "ispresent", ok)
		http.Error(w, "Internal Server", http.StatusInternalServerError)
		return
	}

	docref := firestoreClient.Doc(obj)
	if docref == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	// Get the firestore document
	doc, err := docref.Get(r.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		slog.Error("Firestore doc get error", "error", err)
		http.Error(w, "Internal Server", http.StatusInternalServerError)
		return
	}

	data := doc.Data()

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("Firestore doc encode error", "error", err)
		http.Error(w, "Internal Server", http.StatusInternalServerError)
	}

}
