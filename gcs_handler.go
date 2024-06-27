package gcpwebserv

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	_ "google.golang.org/grpc/balancer/rls"
	_ "google.golang.org/grpc/xds/googledirectpath"
)

// Holding the client as a global since it is concurent safe.
// there is a minor performance gain create a global var.
var gcsClient *storage.Client

// RegisterGCSHandler registers a prebuilt handler for GCS. Custom options can
// be provided using the opts param. A composable http url path is used to
// to get any file in gcs in the form https://domain.com/gcs/bucket/folder/file.file
// In addition, a global reusable storage client is created and maintained in
// this package and will closed when the server shuts down.
func RegisterGCSHandler(ctx context.Context, opts ...option.ClientOption) error {
	// Defaulting the client to use JSON reads. It will become the default
	// in the storage package at some point.
	opts = append(opts, storage.WithJSONReads())
	var err error

	// Creating a GCS client for reuse
	gcsClient, err = storage.NewClient(ctx, opts...)
	if err != nil {
		return err
	}
	clientManager.Add("gcs", gcsClient)

	Route("/gcs/", getObjectFromGCS)
	return nil
}

func getObjectFromGCS(w http.ResponseWriter, r *http.Request) {
	// Excepted path https://domain.com/gcs/bucket/folder/file.file
	// Clean the path of the url
	path := strings.TrimPrefix(path.Clean(r.URL.Path), "/")

	// Split the path and verify min length to make the request
	sPath := strings.Split(path, "/")
	if len(sPath) < 3 {
		http.Error(w, "Bad Path", http.StatusBadRequest)
		return
	}

	// Rejoin the path getting just the file path
	obj := strings.Join(sPath[2:len(sPath)], "/")

	// Attempt to get the file requested
	item, err := gcsClient.Bucket(sPath[1]).Object(obj).NewReader(r.Context())
	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		slog.Error("Error in attempt to get file from bucket", "error", err)
		http.Error(w, "Internal Server", http.StatusInternalServerError)
		return
	}

	// Attemp to read the file
	bts, err := io.ReadAll(item)
	if err != nil {
		slog.Error("Error in attempt read the file", "error", err)
		http.Error(w, "Internal Server", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bts)
	if err != nil {
		slog.Error("Error in attempt read the file", "error", err)
		http.Error(w, "Internal Server", http.StatusInternalServerError)
		return
	}

}
