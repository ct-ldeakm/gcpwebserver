package gcpwebserv

// import (
// 	"bytes"
// 	"context"
// 	"encoding/gob"
// 	"encoding/json"
// 	"fmt"
// 	"log/slog"
// 	"net/http"
// 	"reflect"

// 	"cloud.google.com/go/bigtable"
// 	"google.golang.org/api/option"
// )

// func RegisterBigtableHandler(ctx context.Context, projectId string, instanceId string, opts ...option.ClientOption) error {
// 	client, err := bigtable.NewClient(ctx, projectId, instanceId, opts...)
// 	if err != nil {
// 		slog.Error("Error creating Bigtable Client", "error", err)
// 		return err
// 	}
// 	clientManager.Add(fmt.Sprintf("bigtable-%s-%s", projectId, instanceId), client)

// 	Route("GET /bigtable/{projectId}/{instanceId}/{table}/{key}", btReadRowHandler)
// 	return nil
// }

// func btReadRowHandler(w http.ResponseWriter, r *http.Request) {
// 	slog.Info("In Readrow")
// 	projectId := r.PathValue("projectId")
// 	instanceId := r.PathValue("instanceId")
// 	table := r.PathValue("table")
// 	key := r.PathValue("key") //r.PathValue("key")

// 	slog.Info("Key Requested:", "key", r.URL.Path)

// 	client, ok := clientManager.Get(fmt.Sprintf("bigtable-%s-%s", projectId, instanceId)).(*bigtable.Client)
// 	if !ok {
// 		slog.Error("Unable to retrieve client for the request", "isclient", ok)
// 		http.Error(w, "internal server", http.StatusBadRequest)
// 		return
// 	}

// 	bttable := client.Open(table)
// 	if bttable == nil {
// 		slog.Error("Bigtable table doesnot exits", "table", table)
// 		http.Error(w, "Bigtable table doesnot exits", http.StatusNotFound)
// 		return
// 	}
// 	row, err := client.Open(table).ReadRow(r.Context(), key, bigtable.RowFilter(bigtable.LatestNFilter(1)))
// 	if err != nil {
// 		slog.Error("Bigtable table key doesnot exits", "key", key)
// 		http.Error(w, "Bigtable table key doesnot exits", http.StatusNotFound)
// 		return
// 	}

// 	resultRow := make(map[string]interface{})
// 	for _, v := range row {
// 		for _, item := range v {
// 			dv, err := decodeValue(item.Value)
// 			if err != nil {
// 				slog.Error("Bad Encoded Value", "key", key, "error", err)
// 				http.Error(w, "Internal Server", http.StatusInternalServerError)
// 			}
// 			resultRow[item.Column] = dv
// 		}
// 	}

// 	err = json.NewEncoder(w).Encode(&resultRow)
// 	if err != nil {
// 		slog.Error("Bigtable row encoding error", "err", err)
// 		http.Error(w, "Internal server", http.StatusInternalServerError)
// 	}

// }

// func RegisterBTGetTablesHandler(ctx context.Context, projectId string, instanceId string, opts ...option.ClientOption) error {
// 	client, err := bigtable.NewAdminClient(ctx, projectId, instanceId, opts...)
// 	if err != nil {
// 		slog.Error("Error creating Bigtable Client", "error", err)
// 		return err
// 	}
// 	clientManager.Add(fmt.Sprintf("bigtableadmin-%s-%s", projectId, instanceId), client)

// 	Route("GET /bigtableadmin/{projectId}/{instanceId}/tables", btListTablesHandler)
// 	return nil
// }

// func btListTablesHandler(w http.ResponseWriter, r *http.Request) {

// 	projectId := r.PathValue("projectId")
// 	instanceId := r.PathValue("instanceId")

// 	adminClient, ok := clientManager.Get(fmt.Sprintf("bigtableadmin-%s-%s", projectId, instanceId)).(*bigtable.AdminClient)
// 	if !ok {
// 		slog.Error("Unable to retrieve adminclient for the request", "isclient", ok)
// 		http.Error(w, "internal server", http.StatusBadRequest)
// 		return
// 	}

// 	tables, err := adminClient.Tables(r.Context())
// 	if err != nil {
// 		slog.Error("Unable to get table list", "error", err)
// 		http.Error(w, "internal server", http.StatusBadRequest)
// 		return
// 	}

// 	err = json.NewEncoder(w).Encode(tables)
// 	if !ok {
// 		slog.Error("Unable to encode results", "error", err)
// 		http.Error(w, "internal server", http.StatusBadRequest)
// 		return
// 	}

// }

// func decodeValue(bts []byte) (interface{}, error) {

// 	reflect.TypeOf()
// 	var item interface{}
// 	reader := bytes.NewReader(bts)
// 	err := gob.NewDecoder(reader).Decode(&item)

// 	return item, err

// }
