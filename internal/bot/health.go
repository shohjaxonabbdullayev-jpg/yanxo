package bot

import (
	"context"
	"log"
	"net/http"
	"time"
)

func newHealthMux() http.Handler {
	mux := http.NewServeMux()
	h := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}
	mux.HandleFunc("GET /", h)
	mux.HandleFunc("GET /health", h)
	mux.HandleFunc("GET /healthz", h)
	return mux
}

// startHealthServer listens on addr (e.g. ":8080"). Caller must Shutdown the returned server on shutdown.
func startHealthServer(addr string) *http.Server {
	srv := &http.Server{
		Addr:              addr,
		Handler:           newHealthMux(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		log.Printf("health HTTP server listening on %s (/, /health, /healthz)", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("health HTTP server: %v", err)
		}
	}()
	return srv
}

func shutdownHealthServer(ctx context.Context, srv *http.Server) {
	if srv == nil {
		return
	}
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("health HTTP shutdown: %v", err)
	}
}
