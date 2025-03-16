package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error{
	s.server = &http.Server{
		Addr: ":" + port,
		Handler: handler,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 10,
	}
	 
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error{
	return s.server.Shutdown(ctx)
}