package dao

import (
	"fmt"
	"os"
	"path/filepath"
)

type Image struct {
	Content  []byte
	FileName string
	FilePath string
}

func (i *Image) StoreImage() (bool, error) {
	if _, err := os.Stat(filepath.Join(i.FilePath, i.FileName)); os.IsNotExist(err) {
		// Create file and write to file
		f, err := os.Create(filepath.Join(i.FilePath, i.FileName))
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		defer f.Close()
		if _, err = f.Write(i.Content); err != nil {
			return false, err
		}
	}
	return true, nil
}
