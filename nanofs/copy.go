package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// it can't be named copy due to duplicate name
func copyFiles(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("copy ParseForm error")
		return
	}

	srcPath := r.URL.Path
	dstPath := strings.Join(r.Form["to"], "")
	if !check(srcPath) || !check(dstPath) {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("copy %s to %s: path error", srcPath, dstPath)
		return
	}
	srcPath = "." + srcPath
	dstPath = "." + dstPath

	f, err := os.Stat(srcPath)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		log.Printf("copy %s to %s: %s", srcPath, dstPath, err.Error())
		return
	}

	var result string
	if f.IsDir() {
		result = copyDir(srcPath, dstPath)
	} else {
		result = copyFile(srcPath, dstPath)
	}
	if result != "" {
		http.Error(w, "500", http.StatusInternalServerError)
		log.Printf("copy %s to %s: %s", srcPath, dstPath, result)
	} else {
		_, _ = w.Write([]byte("success"))
		log.Printf("copy %s to %s success", srcPath, dstPath)
	}
}

func copyFile(srcPath, dstPath string) string {
	src, err := os.Open(srcPath)
	if err != nil {
		log.Printf("err : %s", err)
		return err.Error()
	}
	defer src.Close()
	dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("err : %s", err)
		return err.Error()
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		log.Printf("err : %s", err)
		return err.Error()
	}
	return ""
}

func copyDir(src, dst string) string {
	if os.Mkdir(dst, os.ModePerm) != nil {
		return "mkdir error"
	}
	result := ""
	files, err := ioutil.ReadDir(src)
	if err != nil {
		log.Printf("err : %s", err)
		return err.Error()
	}
	for _, f := range files {
		if f.Name() == filepath.Base(src) {
			continue
		}
		if f.IsDir() {
			if result = copyDir(src+"/"+f.Name(), dst+"/"+f.Name()); result != "" {
				break
			}
		} else if result = copyFile(src+"/"+f.Name(), dst+"/"+f.Name()); result != "" {
			break
		}
	}
	return result
}
