package main

import (
	. "io/ioutil"
	"log"
	"net/http"
	"strings"
)

func download(w http.ResponseWriter, r *http.Request) {
	srcPath := r.URL.Path
	if !check(srcPath) {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("download %s: path error", srcPath)
		return
	}

	data, err := ReadFile("." + srcPath)
	if err == nil {
		w.Header().Add("Content-Disposition", "attachment")
		index := strings.LastIndex(srcPath, "/")
		filename := "unknown"
		if index > 0 {
			filename = srcPath[index:]
		}
		w.Header().Add("filename", filename)
		_, err = w.Write(data)
		if err == nil {
			return
		}
	}
	http.Error(w, "404", http.StatusNotFound)
}
