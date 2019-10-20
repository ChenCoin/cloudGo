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
	Port int    `json:"port"`;
	Dir  string `json:"dir"`
}

var defaultConfig = `{
	"port": 80,
	"dir": "."
}`

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
