package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FileType struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"time"`
	IsDir   bool        `json:"isDir"`
}

func list(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !check(path) {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("list %s: path error", path)
		return
	}

	path = "." + path

	files, err := ioutil.ReadDir(path)
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("list.ReadDir %s: %s", path, err.Error())
		return
	}

	var fileData []FileType
	for _, f := range files {
		if f.Name() == filepath.Base(path) {
			continue
		}
		fileData = append(fileData, FileType{f.Name(), f.Size(), f.Mode(), f.ModTime(), f.IsDir()})
	}

	jsonStr, err := json.Marshal(fileData)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		log.Printf("list %s: %s", path, err.Error())
	} else {
		_, _ = w.Write([]byte(string(jsonStr)))
		log.Printf("list %s success", path)
	}
}
