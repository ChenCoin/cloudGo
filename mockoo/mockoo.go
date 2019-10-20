package main

import (
	"encoding/json"
	. "io/ioutil"
	. "log"
	"net/http"
)

var configPath = "./mockoo.json"

var defaultConfig = `{
  "port": ":8090",
  "textRoute": [],
  "fileRoute": []
}`

type MockConfig struct {
	Port       string `json:"port"`
	TextRouter string `json:"textRoute"`
	FileRouter string `json:"fileRoute"`
}

type Router struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RouterList struct {
	Router []Router `json:"router"`
}

func readConfig() (MockConfig, error) {
	data, err := ReadFile(configPath)
	if err != nil {
		return MockConfig{}, err
	}
	conf := MockConfig{}
	err = json.Unmarshal(data, &conf)
	return conf, err
}

func writeConfig() {
	err := WriteFile(configPath, []byte(defaultConfig), 0644)
	if err != nil {
		Printf("config.json created error, %s", err.Error())
	} else {
		Printf("config.json had be created, please check and restart server.")
	}
}

func readRouter(path string) (RouterList, error) {
	data, err := ReadFile(path)
	if err != nil {
		return RouterList{}, err
	}
	conf := RouterList{}
	err = json.Unmarshal(data, &conf)
	return conf, err
}

func main() {
	conf, err := readConfig()
	if err != nil {
		writeConfig()
		return
	}

	routeHandler := func(w http.ResponseWriter, r *http.Request) {
		Printf(r.URL.Path)
		var found = false
		fileRouter, err := readRouter(conf.FileRouter)
		if err == nil {
			for _, v := range fileRouter.Router {
				Printf(v.Key)
				if v.Key == r.URL.Path {
					data, err := ReadFile("." + v.Value)
					if err == nil {
						_, _ = w.Write(data)
					} else {
						http.Error(w, "404", http.StatusNotFound)
					}
					found = true
					break
				}
			}
		}

		textRouter, err := readRouter(conf.TextRouter)
		if err == nil {
			Printf("not nil %d", len(textRouter.Router))
			for _, v := range textRouter.Router {
				Printf(v.Key)
				if v.Key == r.URL.Path {
					_, err = w.Write([]byte(v.Value))
					if err == nil {
						found = true
						break
					}
				}
			}
		} else {
			Printf("%s", err.Error())
		}

		if !found {
			http.Error(w, "400", http.StatusBadRequest)
		}
	}

	Println("server start")
	http.HandleFunc("/", routeHandler)
	err = http.ListenAndServe(conf.Port, nil)
	if err != nil {
		Print("server error: " + err.Error())
	}

	Println("server closed")
}
