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
	Initial bool   `json:"initial"`
	Port    int    `json:"port"`
	Direct  []Pair `json:"direct"`
	File    []Pair `json:"file"`
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
		for _, pair := range conf.Direct {
			if r.RequestURI == pair.Path {
				_, err = w.Write([]byte(pair.Result))
				if err != nil {
					http.Error(w, "500", http.StatusInternalServerError)
				}
				return
			}
		}

		for _, pair := range conf.File {
			if r.RequestURI == pair.Path {
				data, err := ReadFile("." + pair.Result)
				if err == nil {
					_, err = w.Write(data)
					if err != nil {
						http.Error(w, "500", http.StatusInternalServerError)
					}
					return
				}
				http.Error(w, "500", http.StatusInternalServerError)
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
