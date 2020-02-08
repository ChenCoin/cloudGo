package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type site struct {
	Host  string `json:"host"`
	Parse string `json:"parse"`
}

type parseConfig struct {
	Initial bool   `json:"initial"`
	Port    int    `json:"port"`
	Crt     string `json:"crt"`
	Key     string `json:"key"`
	List    []site `json:"list"`
}

var configPath = "./config.json"

func readConfig() (parseConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return parseConfig{}, err
	}
	conf := parseConfig{}
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
		match := false
		for i := 0; i < len(conf.List); i++ {
			site := conf.List[i]
			if r.Host == site.Host {
				uri, err := url.Parse(site.Parse)
				if err != nil {
					return
				}
				httputil.NewSingleHostReverseProxy(uri).ServeHTTP(w, r)
				match = true
				break
			}
		}
		if !match {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	log.Printf("server running on port " + strconv.Itoa(conf.Port))
	if conf.Crt != "" {
		err = http.ListenAndServeTLS(":"+strconv.Itoa(conf.Port), conf.Crt, conf.Key, nil)
	} else {
		err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), nil)
	}
	if err != nil {
		log.Print("server error: " + err.Error())
		return
	}
	log.Print("server closed")
}
