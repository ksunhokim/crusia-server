package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/sunho/crusia-server/api"
	"github.com/sunho/crusia-server/store"
)

const (
	gracefulDuration = 5 * time.Second
)

type Secret struct {
	Version int
	Payload []byte
}

type Server struct {
	version   int
	secrets   []Secret
	stor      store.Store
	apiServer *http.Server
}

func New(version int, stor store.Store, secrets []Secret, addr string) *Server {
	s := &Server{
		version: version,
		secrets: secrets,
		stor:    stor,
	}

	a := api.New(&apiInterface{s})
	h := &http.Server{
		Handler: a.Http(),
		Addr:    addr,
	}
	s.apiServer = h

	return s
}

func (s *Server) Run() {
	go s.runApi()
}

func (s *Server) runApi() {
	err := s.apiServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (s *Server) Stop() {
	s.stopApi()
}

func (s *Server) stopApi() {
	ctx, cancel := context.WithTimeout(context.Background(), gracefulDuration)
	defer cancel()

	if err := s.apiServer.Shutdown(ctx); err != nil {
		log.Println("ERROR", "failed to shutdown server", err)
	} else {
		log.Println("INFO", "server stopped")
	}
}
