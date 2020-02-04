package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func move(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("move ParseForm: path error")
		return
	}

	srcPath := r.URL.Path
	dstPath := strings.Join(r.Form["to"], "")

	if !check(srcPath) || !check(dstPath) {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("move %s %s: path error", srcPath, dstPath)
		return
	}

	err = os.Rename("."+srcPath, "."+dstPath)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		log.Printf("move %s %s: %s", srcPath, dstPath, err.Error())
	} else {
		_, _ = w.Write([]byte("success"))
		log.Printf("move %s %s success", srcPath, dstPath)
	}
}
