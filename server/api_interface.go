package server

import "github.com/sunho/crusia-server/store"

type apiInterface struct {
	s *Server
}

func (a *apiInterface) DecryptSaveData(version int, payload string) ([]byte, error) {
	return a.s.decryptSaveData(version, payload)
}

func (a *apiInterface) CreateToken(id int) (string, error) {
	return a.s.createToken(id)
}

func (a *apiInterface) GetToken(tok string) (int, error) {
	return a.s.getToken(tok)
}

func (a *apiInterface) GetStore() store.Store {
	return a.s.stor
}

func (a *apiInterface) GetVersion() int {
	return a.s.version
}
