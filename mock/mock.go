package main

import (
	"encoding/json"
	. "io/ioutil"
	. "log"
	"net/http"
	"sort"
	"strconv"
)

var configPath = "./config.json"

type Pair struct {
	Path   string `json:"path"`
	Result string `json:"result"`
}

type Config struct {
	Initial  bool     `json:"initial"`
	Port     int      `json:"port"`
	Direct   []Pair   `json:"direct"`
	File     []Pair   `json:"file"`
	Redirect []Pair   `json:"redirect"`
	Dir      []string `json:"dir"`
	Extra    []string `json:"extra"`
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

	readExtraConfig(conf)
	sort.Sort(PairSlice(conf.Direct))
	sort.Sort(PairSlice(conf.File))
	sort.Sort(PairSlice(conf.Redirect))
	sort.Sort(StringLenSlice(conf.Dir))

	handler := func(w http.ResponseWriter, r *http.Request) {
		// uri := r.URL.Path is another resolution
		// if the request is http://localhost:8090/money?name=andy
		// then r.URL.Path will be /money
		// and r.RequestURI will be /money?name=andy
		// it is better to match /money?name=andy
		uri := r.RequestURI
		path := r.URL.Path
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

		for _, pair := range conf.File {
			if uri == pair.Path {
				data, err := ReadFile(pair.Result)
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

		for _, pair := range conf.Redirect {
			if uri == pair.Path {
				data, err := ReadFile(pair.Result)
				if err == nil {
					redirectData, err := ReadFile(string(data))
					if err == nil {
						_, err = w.Write(redirectData)
						if err == nil {
							return
						}
					}
				}
				http.Error(w, "404", http.StatusNotFound)
				return
			}
		}

		for _, dir := range conf.Dir {
			if len(path) >= len(dir) && path[0:len(dir)] == dir {
				data, err := ReadFile(path)
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

func readExtraConfig(conf Config) {
	for _, path := range conf.Extra {
		data, err := ReadFile(path)
		if err != nil {
			continue
		}
		extra := Config{}
		err = json.Unmarshal(data, &extra)
		if err != nil {
			continue
		}
		if len(extra.Direct) > 0 {
			conf.Direct = append(conf.Direct, extra.Direct...)
		}
		if len(extra.File) > 0 {
			conf.File = append(conf.File, extra.File...)
		}
		if len(extra.Redirect) > 0 {
			conf.Redirect = append(conf.Redirect, extra.Redirect...)
		}
		if len(extra.Dir) > 0 {
			conf.Dir = append(conf.Dir, extra.Dir...)
		}
	}

}

// sort
type StringLenSlice []string

func (p StringLenSlice) Len() int {
	return len(p)
}

func (p StringLenSlice) Less(i, j int) bool {
	return len(p[i]) < len(p[j])
}

func (p StringLenSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type PairSlice []Pair

func (p PairSlice) Len() int {
	return len(p)
}

func (p PairSlice) Less(i, j int) bool {
	return len(p[i].Path) < len(p[j].Path)
}

func (p PairSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
