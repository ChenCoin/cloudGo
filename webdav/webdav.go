package main

import (
	"encoding/json"
	"golang.org/x/net/webdav"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type WebdavConfig struct {
	Initial bool   `json:"initial"`
	Port    int    `json:"port"`
	Dir     string `json:"dir"`
}

var configPath = "./webdav.json"

func readConfig() (WebdavConfig, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return WebdavConfig{}, err
	}
	conf := WebdavConfig{}
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

	log.Print("server start")
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
