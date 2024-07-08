package gcpwebserv

import (
	"context"
	"fmt"

	"io"
	"log/slog"
	"net/http"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func RegisterPubSubHandler(ctx context.Context, projectId string, topicId string, opts ...option.ClientOption) error {
	client, err := pubsub.NewClient(ctx, projectId, opts...)
	if err != nil {
		slog.Error("Error building pubsub client", "error", err)
		return err
	}

	AddClientToCache(fmt.Sprintf("pubsub-%s-%s", projectId, topicId), client)

	Route("POST /pubsub/{projectId}/{topicId}", publishPubSubMessage)
	return nil
}

func publishPubSubMessage(w http.ResponseWriter, r *http.Request) {
	bdy, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		slog.Error("Request body read error", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if bdy == nil {
		slog.Error("Empty body", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	projectId := r.PathValue("projectId")
	topicId := r.PathValue("topicId")

	client, ok := GetCachedClient(fmt.Sprintf("pubsub-%s-%s", projectId, topicId)).(*pubsub.Client)
	if !ok {
		slog.Error("Unable to retreive client from cache", "ispresent", ok)
		http.Error(w, "Internal Server", http.StatusInternalServerError)
		return
	}
	topic := client.TopicInProject(topicId, projectId)

	result := topic.Publish(r.Context(), &pubsub.Message{
		Data: bdy,
	})

	_, err = result.Get(r.Context())
	if err != nil {
		slog.Error("Error Publishing pubsub message", "error", err)
		http.Error(w, "Internal Server", http.StatusInternalServerError)
	}

}
