package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	defaultAddr     = ":8080"
	shutdownTimeout = 15 * time.Second
)

var (
	ErrServerAlreadyStarted = errors.New("server already started")
)

type Server struct {
	httpServer *http.Server
	isStarted  bool
}

func NewServer(handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         defaultAddr,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,  // Time to read the request
			WriteTimeout: 15 * time.Second, // Response time (including nsjail)
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	if s.isStarted {
		return ErrServerAlreadyStarted
	}

	s.isStarted = true
	defer func() {
		s.isStarted = false
	}()

	go func() {
		<-ctx.Done()
		log.Println("🛑 Shutting down worker gracefully...")
		err := s.shutdownWithTimeout(shutdownTimeout)
		if err != nil {
			log.Printf("Failed to shutdown server: %v", err)
		}
	}()

	port := strings.TrimPrefix(s.httpServer.Addr, ":")
	log.Printf("🏰 Kube-Sandcastle HTTP Worker starting on %s...", port)
	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}

func (s *Server) shutdownWithTimeout(timeout time.Duration) error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.httpServer.Shutdown(shutdownCtx)
}
