package server

import (
	"log"
	"net/http"

	"github.com/sunho/crusia-server/store"
)

type Secret struct {
	Version int
	Payload []byte
}

type Server struct {
	Version    int
	Secrets    []Secret
	Store      store.Store
	HttpServer *http.Server
}

func New(version int, stor store.Store, secrets []Secret, addr string) *Server {
	s := &Server{
		Version: version,
		Secrets: secrets,
		Store:   stor,
	}

	h := &http.Server{
		Addr: addr,
	}
	s.HttpServer = h
	return s
}

func (s *Server) Run() {
	go s.runApi()
}

func (s *Server) runApi() {
	err := s.HttpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
