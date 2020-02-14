package main

import (
	"encoding/json"
	. "io/ioutil"
	. "log"
	"net/http"
	"strconv"
)

var configPath = "./config.json"

type Pair struct {
	Path   string `json:"path"`
	Result string `json:"result"`
}

type Config struct {
	Initial bool     `json:"initial"`
	Port    int      `json:"port"`
	Direct  []Pair   `json:"direct"`
	Dir     []string `json:"dir"`
	File    []Pair   `json:"file"`
}

func readConfig() (Config, error) {
	data, err := ReadFile(configPath)
	conf := Config{}
	if err == nil {
		err = json.Unmarshal(data, &conf)
	}
	return conf, err
}

func main() {
	conf, err := readConfig()
	if err != nil {
		Print(err.Error())
		return
	}

	if conf.Initial == false {
		Print("please initial the configure file, and set \"initial\":true")
		return
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		// uri := r.URL.Path
		// if the request is http://localhost:8090/money?name=andy
		// then r.URL.Path will be /money
		// and r.RequestURI will be /money?name=andy
		uri := r.RequestURI
		Print("please initial the config " + uri + " " + r.URL.Path)
		for _, pair := range conf.Direct {
			if uri == pair.Path {
				_, err = w.Write([]byte(pair.Result))
				if err != nil {
					http.Error(w, "500", http.StatusInternalServerError)
				}
				return
			}
		}

		for _, dir := range conf.Dir {
			if len(uri) >= len(dir) && uri[0:len(dir)] == dir {
				data, err := ReadFile("." + uri)
				if err == nil {
					_, err = w.Write(data)
					if err != nil {
						http.Error(w, "500", http.StatusInternalServerError)
					}
					return
				}
				http.Error(w, "404", http.StatusNotFound)
				return
			}
		}

		for _, pair := range conf.File {
			if uri == pair.Path {
				data, err := ReadFile("." + pair.Result)
				if err == nil {
					_, err = w.Write(data)
					if err != nil {
						http.Error(w, "500", http.StatusInternalServerError)
					}
					return
				}
				http.Error(w, "404", http.StatusNotFound)
				return
			}
		}
		http.Error(w, "404", http.StatusNotFound)
	}

	Printf("server running on port " + strconv.Itoa(conf.Port))
	http.HandleFunc("/", handler)
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), nil)
	if err != nil {
		Print("server error: " + err.Error())
	}
	Println("server closed")
}
