package server

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sdbx/crusia-server/api"
	"github.com/sdbx/crusia-server/store"
	"github.com/sdbx/crusia-server/utils"
)

const (
	gracefulDuration = 5 * time.Second
	noiseSize        = 20
)

var (
	ErrInvalidVersion  = errors.New("server: invalid version")
	ErrInvalidSaveData = errors.New("server: save data is not in json format")
)

type SaveKey struct {
	Version int
	IV      []byte
	Payload []byte
}

type Server struct {
	version   int
	key       []byte
	iv        []byte,
	saveKeys  []SaveKey
	stor      store.Store
	apiServer *http.Server
}

func New(version int, stor store.Store, key []byte, iv []byte, saveKeys []SaveKey, addr string) *Server {
	s := &Server{
		version:  version,
		key:      key,
		iv:       iv,
		saveKeys: saveKeys,
		stor:     stor,
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
	log.Println("INFO", "server started", s.apiServer.Addr)
	err := s.apiServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (s *Server) createToken(id int) (string, error) {
	str := strconv.Itoa(id)
	payload := make([]byte, noiseSize, noiseSize+len(str))
	if _, err := io.ReadFull(rand.Reader, payload); err != nil {
		return "", err
	}

	payload = append(payload, []byte(str)...)
	return Encrypt(s.key, s.iv, payload)
}

func (s *Server) getToken(tok string) (int, error) {
	payload, err := Decrypt(s.key, s.iv, tok)
	if err != nil {
		return 0, err
	}

	payload = payload[noiseSize:]
	id, err := strconv.Atoi(string(payload))
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Server) getSaveKey(version int) (SaveKey, bool) {
	for _, key := range s.saveKeys {
		if key.Version == version {
			return key, true
		}
	}

	return nil, false
}

func (s *Server) decryptSaveData(version int, payload string) ([]byte, error) {
	key, ok := s.getSaveKey(version)
	if !ok {
		return nil, ErrInvalidVersion
	}

	buf, err := Decrypt(key.Payload, key.IV, payload)
	if err != nil {
		return nil, err
	}
	log.Println(string(buf))
	if !utils.IsJSON(buf) {
		return nil, ErrInvalidSaveData
	}

	return buf, nil
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
