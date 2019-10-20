package main

import (
	"golang.org/x/net/webdav"
	"net/http"
)

func main() {
	_ = http.ListenAndServe(":8848", &webdav.Handler{
		FileSystem: webdav.Dir("."),
		LockSystem: webdav.NewMemLS(),
	})
}
