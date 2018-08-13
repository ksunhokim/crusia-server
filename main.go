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

type Secret struct {
	Version int    `yaml:"version"`
	Key     string `yaml:"key"`
}

type Config struct {
	Addr    string   `yaml:"addr"`
	Version int      `yaml:"version"`
	Secrets []Secret `yaml:"secrets"`
}

func parseSecrets(secrets []Secret) ([]server.Secret, error) {
	out := make([]server.Secret, len(secrets))
	for _, s := range secrets {
		buf, err := base64.StdEncoding.DecodeString(s.Key)
		if err != nil {
			return nil, err
		}
		out = append(out, server.Secret{
			Version: s.Version,
			Payload: buf,
		})
	}
	return out, nil
}

func main() {
	buf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	conf := Config{}
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		log.Fatal(err)
	}

	secrets, err := parseSecrets(conf.Secrets)
	if err != nil {
		log.Fatal(err)
	}

	db, err := storm.Open("bolt.db")
	if err != nil {
		log.Fatal(err)
	}

	bs := boltstore.New(db)
	err = bs.Init()
	if err != nil {
		log.Fatal(err)
	}

	serv := server.New(conf.Version, bs, secrets, conf.Addr)
	serv.Run()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	serv.Stop()
	db.Close()
}
