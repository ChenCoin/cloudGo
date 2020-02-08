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

	http.Handle("/", http.FileServer(http.Dir(conf.Dir)))
	log.Printf("server running on " + strconv.Itoa(conf.Port))
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), nil)
	if err != nil {
		log.Printf("error when create server")
	}
}
