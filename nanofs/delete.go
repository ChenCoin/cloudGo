package main

import (
	"log"
	"net/http"
	"os"
)

// it can't be named delete due to duplicate name
func deleteFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !check(path) {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("delete %s: path error", path)
		return
	}
	err := os.RemoveAll("." + path)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		log.Printf("delete %s: %s", path, err.Error())
	} else {
		_, _ = w.Write([]byte("success"))
		log.Printf("delete %s success", path)
	}
}
