package service

import (
	"path"
	"strings"

	"github.com/sysu-activitypluspc/service-end/dao"
)

type Image struct {
	dao.Image
	FileName string
}

func (i *Image) UploadImage() (int, error) {
	md5Filename := GetMd5(i.Content)
	ext := path.Ext(i.FileName)
	filename := strings.Join([]string{md5Filename, ext}, "")
	i.Image.FileName = filename
	_, err := i.StoreImage()
	if err != nil {
		return 500, err
	}
	return 200, nil
}
