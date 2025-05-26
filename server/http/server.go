package server

import (
	"context"
	"net/http"
	"time"

	"github.com/b85bagent/tools/server"
)

type HTTPServer struct {
	srv *http.Server
}

func NewHTTPServer(addr string, handler http.Handler) server.Server {
	return &HTTPServer{
		srv: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (h *HTTPServer) Start() error {
	return h.srv.ListenAndServe()
}

func (h *HTTPServer) Stop(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return h.srv.Shutdown(shutdownCtx)
}
