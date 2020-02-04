package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ParseConfig struct {
	Initial bool   `json:"initial"`
	Address string `json:"address"`
	Crt     string `json:"crt"`
	Key     string `json:"key"`
	Parse   string `json:"parse"`
}

var configPath = "./proxyGo.json"

func readConfig() (ParseConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return ParseConfig{}, err
	}
	conf := ParseConfig{}
	err = json.Unmarshal(data, &conf)
	return conf, err
}

func main() {
	conf, err := readConfig()
	if err != nil {
		log.Print(err.Error())
		return
	}

	if conf.Initial == false {
		log.Print("please initial the configure file, and set \"initial\":true")
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uri, err := url.Parse(conf.Parse)
		if err != nil {
			return
		}
		httputil.NewSingleHostReverseProxy(uri).ServeHTTP(w, r)
	})
	log.Print("server start")
	err = http.ListenAndServeTLS(conf.Address, conf.Crt, conf.Key, nil)
	if err != nil {
		log.Print("server error: " + err.Error())
		return
	}
	log.Print("server closed")
}
