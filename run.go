package gcpwebserv

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run starts the ListenAndServe server in a go routine with the global
// mux as the handler. The function blocks waiting on a shutddown of
// either Sig Term (use by GCP service) or Sig Int (Ctrl+c).
func Run(ctx context.Context, server *http.Server) error {

	server.Handler = mux
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
		log.Println("Rejecting new connections.")
	}()

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, 15*time.Second)
	defer shutdownRelease()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		var sd bool
		select {
		case sig := <-sigChan:
			sd = true
			slog.Info("Shutdown signal received", "signal", sig)
		case <-ctx.Done():
			sd = true
			slog.Info("Context done signaled")
		}
		if sd {
			break
		}

	}

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		slog.Error("Server Shutdown Error", "error", err)
	}

	// Closing all clients for gcs services
	closeAllClients()
	slog.Info("Server Exiting", "status", err)
	return err
}
