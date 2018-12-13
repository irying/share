package upload

import (
	"backend/conf"
	"backend/utils"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// CheckImageExt check extension of image
func CheckImageExt(fileName string) bool {
	ext := GetExt(fileName)
	for _, allowExt := range conf.Conf.App.Api.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

// CheckImageSize check size of image
func CheckImageSize(f multipart.File) bool {
	size, err := GetSize(f)
	if err != nil {
		utils.LogPrintError(err)
		return false
	}

	return size <= conf.Conf.App.Api.ImageMaxSize
}

// GetImageName get name of image
func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = utils.EncodeMD5(fileName)

	return fileName + ext
}

// GetImagePath get path of image
func GetImagePath() string {
	dir, err := os.Getwd()
	utils.LogPrintError(err)

	return dir + "/../" + conf.Conf.App.Api.ImageSavePath
}