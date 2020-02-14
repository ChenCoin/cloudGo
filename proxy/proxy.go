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

type Rewrite struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Site struct {
	Host    string  `json:"host"`
	Prefix  string  `json:"Prefix"`
	Proxy   string  `json:"proxy"`
	Rewrite Rewrite `json:"rewrite"`
}

type Config struct {
	Initial bool   `json:"initial"`
	Port    int    `json:"port"`
	Crt     string `json:"crt"`
	Key     string `json:"key"`
	List    []Site `json:"list"`
}

var configPath = "./config.json"

func readConfig() (Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}
	conf := Config{}
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
			match = reverseProxy(site, w, r)
			if match {
				break
			}
		}
		if !match {
			log.Println("no")
			http.Error(w, "404", http.StatusNotFound)
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

func reverseProxy(site Site, w http.ResponseWriter, r *http.Request) bool {
	if site.Host != r.Host {
		return false
	}
	requestURI := r.RequestURI
	if len(requestURI) >= len(site.Prefix) && requestURI[0:len(site.Prefix)] == site.Prefix {
		uri, err := url.Parse(site.Proxy)
		if err != nil {
			log.Print("server error: " + err.Error())
			return false
		}
		rewrite := site.Rewrite
		if rewrite.From != "" || rewrite.To != "" {
			if len(requestURI) >= len(rewrite.From) && requestURI[0:len(rewrite.From)] == rewrite.From {
				r.URL.Path = rewrite.To + requestURI[len(rewrite.From):]
			} else {
				log.Print("the rewrite of [" + site.Host + "," + site.Prefix +
					"] is error that the from can not match prefix")
			}
		}
		httputil.NewSingleHostReverseProxy(uri).ServeHTTP(w, r)
		return true
	}
	return false
}
