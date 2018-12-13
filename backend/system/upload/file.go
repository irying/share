package upload

import (
	"io/ioutil"
	"mime/multipart"
	"path"
)

// GetSize return size of file
func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content), err
}

// GetExt  return extension of file
func GetExt(fileName string) string {
	return path.Ext(fileName)
}