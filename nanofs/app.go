package main

import (
	"log"
	"net/http"
)

func main() {
	handlerFunc("/list/", "/list", list)
	handlerFile("/download/", "/download", "./")
	http.HandleFunc("/put/file/", uploadFile)
	handlerFunc("/upload/", "/upload", uploadFile)
	handlerFunc("/copy/", "/copy", copyFiles)
	handlerFunc("/move/", "/move", move)
	handlerFunc("/delete/", "/delete", deleteFile)
	handlerFunc("/rename/", "/rename", rename)
	handlerFunc("/mkdir/", "/mkdir", mkdir)
	handlerFile("/", "", "./")
	log.Printf("server running")
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Printf("error when create server")
	}
}

func handlerFunc(pattern string, prefix string, handlerFunc http.HandlerFunc) {
	http.Handle(pattern, http.StripPrefix(prefix, handlerFunc))
}

func handlerFile(pattern string, prefix string, path string) {
	http.Handle(pattern, http.StripPrefix(prefix, http.FileServer(http.Dir(path))))
}
