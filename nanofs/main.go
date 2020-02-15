package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Config struct {
	Initial bool   `json:"initial"`
	Port    int    `json:"port"`
	Dir     string `json:"dir"`
}

var configPath = "./config.json"

var root = "."

func readConfig() (Config, error) {
	data, err := ioutil.ReadFile(configPath)
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

	root = conf.Dir
	handlerFunc("/list/", "/list", list)
	handlerFunc("/download/", "/download", download)
	handlerFunc("/upload/", "/upload", uploadFile)
	handlerFunc("/copy/", "/copy", copyFiles)
	handlerFunc("/move/", "/move", move)
	handlerFunc("/delete/", "/delete", deleteFile)
	handlerFunc("/rename/", "/rename", rename)
	handlerFunc("/mkdir/", "/mkdir", mkdir)
	handlerFile("/", "", root+"/")
	log.Printf("server running on port " + strconv.Itoa(conf.Port))
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), nil)
	if err != nil {
		log.Printf("error when create server")
	}
}

func handlerFunc(pattern string, prefix string, handlerFunc http.HandlerFunc) {
	http.Handle(pattern, http.StripPrefix(prefix, handlerFunc))
}

func handlerFile(pattern string, prefix string, path string) {
	http.Handle(pattern, http.StripPrefix(prefix, http.FileServer(http.Dir(path))))
}
