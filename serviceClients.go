package gcpwebserv

import (
	"log/slog"
	"sync"
)

var clientManager serviceClientCache

func init() {
	clientManager.clientList = make(map[string]serviceClient)
}

type serviceClient interface {
	Close() error
}
type serviceClientCache struct {
	mu         sync.Mutex
	clientList map[string]serviceClient
}

func (s *serviceClientCache) Add(name string, client serviceClient) {
	s.mu.Lock()
	s.clientList[name] = client
	s.mu.Unlock()
}
func (s *serviceClientCache) Get(name string) serviceClient {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clientList[name]
}

// CloseAll is used to close all references to clients that are cached.
func (s *serviceClientCache) CloseAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for name, client := range s.clientList {
		err := client.Close()
		if err != nil {
			slog.Error("Error Closing Client", "client", name, "Error", err)
			continue
		}
		slog.Info("Closed Client", "client", name)
	}
}
