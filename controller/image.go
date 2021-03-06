package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sysu-activitypluspc/service-end/types"
	"github.com/sysu-activitypluspc/service-end/services"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	var maxMemory int64 = 5 * (1 << 20)
	r.ParseMultipartForm(maxMemory)
	file, handler, err := r.FormFile("file")
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
	md5Filename := services.GetMd5(content)
	ext := path.Ext(handler.Filename)
	filename := strings.Join([]string{md5Filename, ext}, "")
	// Check if the file exists
	if _, err = os.Stat(filepath.Join(staticFilePosition, filename)); os.IsNotExist(err) {
		// Create file and write to file
		f, err := os.Create(filepath.Join(staticFilePosition, filename))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
		}
		defer f.Close()
		if _, err = f.Write(content); err != nil {
			w.WriteHeader(500)
			return
		}
	}
	fileInfo := types.FileInfo{
		Filename: filename,
	}
	resBody, _ := json.Marshal(fileInfo)
	w.Write(resBody)
}