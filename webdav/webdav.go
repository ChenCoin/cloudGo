package main

import (
	"encoding/json"
	"golang.org/x/net/webdav"
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

	log.Printf("server running on port " + strconv.Itoa(conf.Port))
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), &webdav.Handler{
		FileSystem: webdav.Dir(conf.Dir),
		LockSystem: webdav.NewMemLS(),
	})
	if err != nil {
		log.Print("server error: " + err.Error())
		return
	}
	log.Print("server closed")
}
