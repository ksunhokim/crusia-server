package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/asdine/storm"
	"github.com/sunho/crusia-server/server"
	"github.com/sunho/crusia-server/store/boltstore"
	yaml "gopkg.in/yaml.v2"
)

type SaveKey struct {
	Version int    `yaml:"version"`
	Key     string `yaml:"key"`
}

type Config struct {
	Addr     string    `yaml:"addr"`
	Version  int       `yaml:"version"`
	Key      string    `yaml:"key"`
	SaveKeys []SaveKey `yaml:"save_keys"`
}

func parseSaveKeys(keys []SaveKey) ([]server.SaveKey, error) {
	out := make([]server.SaveKey, len(keys))
	for _, s := range keys {
		buf, err := base64.StdEncoding.DecodeString(s.Key)
		if err != nil {
			return nil, err
		}
		out = append(out, server.SaveKey{
			Version: s.Version,
			Payload: buf,
		})
	}
	return out, nil
}

func createServer(buf []byte, db *storm.DB) (*server.Server, error) {
	conf := Config{}
	err := yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, err
	}

	key, err := base64.StdEncoding.DecodeString(conf.Key)
	if err != nil {
		return nil, err
	}

	skeys, err := parseSaveKeys(conf.SaveKeys)
	if err != nil {
		return nil, err
	}

	bs := boltstore.New(db)
	err = bs.Init()
	if err != nil {
		return nil, err
	}

	serv := server.New(conf.Version, bs, key, skeys, conf.Addr)
	return serv, nil
}

func main() {
	buf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	db, err := storm.Open("bolt.db")
	if err != nil {
		log.Fatal(err)
	}

	serv, err := createServer(buf, db)
	if err != nil {
		log.Fatal(err)
	}

	serv.Run()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	serv.Stop()
	db.Close()
}
