package main

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	page := `
<html>
<head>
    <title>upload</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0, user-scalable=no">
</head>
<body>
	<form enctype="multipart/form-data" action="#" method="post">
		<input type="file" name="files" />
		<input type="submit" value="upload" />
	</form>
</body>
</html>
`
	if r.Method == http.MethodGet {
		_, _ = io.WriteString(w, page)
		log.Printf("upload: get page")
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("upload: method error")
		return
	}

	path := r.URL.Path
	if !check(path) {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("upload: path invalid")
		return
	}

	err := r.ParseMultipartForm(1024000)
	if err != nil {
		http.Error(w, "404", http.StatusNotFound)
		log.Printf("upload: %s", err.Error())
		return
	}

	files := r.MultipartForm.File["files"]
	result := true
	fileNames := ""
	for i, _ := range files {
		file := files[i]
		fileNames += " " + file.Filename
		err = _saveFile(w, file, path)
		if err != nil {
			result = false
			http.Error(w, "404", http.StatusNotFound)
			log.Printf("upload: %s", err.Error())
			break
		}
	}
	if result {
		_, _ = w.Write([]byte("success"))
		log.Printf("upload: success," + fileNames)
	}
}

func _saveFile(w http.ResponseWriter, fileHeader *multipart.FileHeader, path string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	dst, err := os.Create("." + path + "/" + fileHeader.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return err
	}
	return nil
}
