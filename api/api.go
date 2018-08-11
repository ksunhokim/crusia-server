package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sunho/crusia-server/store"
	"github.com/sunho/crusia-server/utils"
)

type ApiInterface interface {
	Decrypt(buf []byte, version int) ([]byte, error)
	CreateToken(id int) (string, error)
	GetToken(tok string) (int, error)
	GetStore() store.Store
	GetVersion() int
}

type Api struct {
	in ApiInterface
}

func New(in ApiInterface) *Api {
	return &Api{in: in}
}

func (a *Api) GetVersion(w http.ResponseWriter, r *http.Request) {
	v := a.in.GetVersion()
	resp := map[string]interface{}{
		"version": v,
	}
	utils.HttpJson(w, resp)
}

func (a *Api) Login(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Username string `json:"username"`
		Passhash string `json:"passhash"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.HttpError(w, err, 400)
		return
	}

	u, err := a.in.GetStore().GetUserByUsername(req.Username)
	if err != nil {
		utils.HttpError(w, err, 404)
		return
	}

	if u.Passhash != req.Passhash {
		utils.HttpError(w, err, 403)
		return
	}

	tok, err := a.in.CreateToken(u.ID)
	if err != nil {
		utils.HttpError(w, err, 500)
		return
	}

	resp := map[string]interface{}{
		"token": tok,
	}
	utils.HttpJson(w, resp)
}

func (a *Api) Register(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Username string `json:"username"`
		Passhash string `json:"passhash"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.HttpError(w, err, 400)
		return
	}

	u, err := a.in.GetStore().CreateUser(&store.User{
		Username: req.Username,
		Passhash: req.Passhash,
	})
	if err != nil {
		utils.HttpError(w, err, 409)
		return
	}
	_, err = a.in.GetStore().CreateSaveData(&store.SaveData{
		UserID:  u.ID,
		Edited:  time.Now(),
		Payload: "{}",
	})
	if err != nil {
		utils.HttpError(w, err, 500)
		return
	}

	w.WriteHeader(200)
}

func getUser(r *http.Request) *store.User {
	return nil
}

func (a *Api) GetSaveData(w http.ResponseWriter, r *http.Request) {
	u := getUser(r)
	data, err := a.in.GetStore().GetSaveData(u.ID)
	if err != nil {
		utils.HttpError(w, err, 500)
		return
	}
	fmt.Fprintln(w, data.Payload)
}

func (a *Api) PostSaveData(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.HttpError(w, err, 400)
	}

	u := getUser(r)
	data := &store.SaveData{
		UserID:  u.ID,
		Edited:  time.Now(),
		Payload: string(buf),
	}

	err = a.in.GetStore().UpdateSaveData(data)
	if err != nil {
		utils.HttpError(w, err, 500)
		return
	}
	w.WriteHeader(201)
}
