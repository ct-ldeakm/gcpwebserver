package gcpwebserv

import (
	"log/slog"
	"sync"
)

// clientManager is a global variable in this package
// the clients that match this interface are goroutine
// safe.
var clientManager ServiceClientCache

// AddClientToCache add the client to cache of clients. On the first
// call to add the cache is created. All clients in the cache are
// closed when the server is shutdown. By naming the clients uniquely
// multiply client to the same service across different services
// can be maintained. The client can also have different authentication
// schemes. This cache is safe for concurrent use.
func AddClientToCache(name string, client serviceClient) {
	clientManager.add(name, client)
}

// GetCachedClient looks up the cached client from the cache
// The client is returned as an interface and will need to
// be inferred by the caller before use. For example, if a big
// query client is stored name bqclient it can be retrieved and
// used as gcpwebserver.Get("bqclient").(*bigquery.Client)
func GetCachedClient(name string) serviceClient {
	return clientManager.get(name)
}

func closeAllClients() {
	clientManager.closeAll()
}

type serviceClient interface {
	Close() error
}
type ServiceClientCache struct {
	mu         sync.Mutex
	clientList map[string]serviceClient
}

func (s *ServiceClientCache) add(name string, client serviceClient) {
	if s.clientList == nil {
		s.clientList = make(map[string]serviceClient)
	}
	s.mu.Lock()
	s.clientList[name] = client
	s.mu.Unlock()
}
func (s *ServiceClientCache) get(name string) serviceClient {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clientList[name]
}

// CloseAll is used to close all references to clients that are cached.
func (s *ServiceClientCache) closeAll() {
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
