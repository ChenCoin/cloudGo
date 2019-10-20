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
	Address string `json:"address"`;
	Crt     string `json:"crt"`;
	Key     string `json:"key"`;
	Parse   string `json:"parse"`
}

var defaultConfig = `{
	"address": ":443",
	"crt": "server.crt",
	"key": "server.key",
	"parse": "http://localhost:80"
}`

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

func writeConfig() {
	err := ioutil.WriteFile(configPath, []byte(defaultConfig), 0644)
	if err != nil {
		log.Printf("config.json created error, %s", err.Error())
	} else {
		log.Printf("config.json had be created, please check and restart server.")
	}
}

func main() {
	conf, err := readConfig()
	if err != nil {
		writeConfig()
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
