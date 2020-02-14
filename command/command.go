package main

import (
	"encoding/json"
	. "io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

type Item struct {
	Path    string `json:"path"`
	Command string `json:"command"`
}

type Config struct {
	Initial bool   `json:"initial"`
	Port    int    `json:"port"`
	Bash    string `json:"bash"`
	List    []Item `json:"list"`
}

var configPath = "./config.json"

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
		log.Print(err.Error())
		return
	}

	if conf.Initial == false {
		log.Print("please initial the configure file, and set \"initial\":true")
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		match := false
		for _, item := range conf.List {
			if r.RequestURI == "/"+item.Path {
				cmd := exec.Command(conf.Bash, "-c", item.Command)
				bytes, err := cmd.Output()
				if err != nil {
					w.Write([]byte(err.Error()))
				} else {
					w.Write(bytes)
				}
				match = true
				break
			}
		}
		if !match {
			http.Error(w, "404", http.StatusNotFound)
		}
	})
	log.Printf("server running on port " + strconv.Itoa(conf.Port))
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), nil)
	if err != nil {
		log.Printf("error when create server")
	}
}
