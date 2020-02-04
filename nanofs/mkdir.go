package main

import (
	"log"
	"net/http"
	"os"
)

func mkdir(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !check(path) {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("mkdir %s: path error", path)
		return
	}
	err := os.Mkdir("."+path, os.ModePerm)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		log.Printf("mkdir %s: %s", path, err.Error())
	} else {
		_, _ = w.Write([]byte("success"))
		log.Printf("mkdir %s success", path)
	}
}
