package ui

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	var maxMemory int64 = 5 * (1 << 20)
	r.ParseMultipartForm(maxMemory)
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
	}
	defer file.Close()
	staticFilePosition := os.Getenv("STATIC_DIR")
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
}
