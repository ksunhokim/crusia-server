package server

import "github.com/sunho/crusia-server/store"

type apiInterface struct {
	s *Server
}

func (a *apiInterface) Decrypt(buf []byte, version int) ([]byte, error) {
	return nil, nil
}

func (a *apiInterface) CreateToken(id int) (string, error) {
	return "", nil
}

func (a *apiInterface) GetToken(tok string) (int, error) {
	return 0, nil
}

func (a *apiInterface) GetStore() store.Store {
	return a.s.stor
}

func (a *apiInterface) GetVersion() int {
	return a.s.version
}
