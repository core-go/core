package io

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func SaveFiles(rootDirectory string, files []File) error {
	for _, v := range files {
		fullPath := rootDirectory + string(os.PathSeparator) + v.Name
		err := Save(fullPath, v.Content)
		if err != nil {
			return err
		}
	}
	return nil
}

func Save(fullName string, content string) error {
	err := os.MkdirAll(filepath.Dir(fullName), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullName, []byte(content), os.ModePerm)
}

func IsValidPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
