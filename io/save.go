package io

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func Save(fullName string, content string) error {
	err := os.MkdirAll(filepath.Dir(fullName), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullName, []byte(content), os.ModePerm)
}
func List(path string) ([]string, error) {
	var names []string
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	folder, err := ioutil.ReadDir(absPath)
	if err != nil {
		return names, err
	}
	for _, tmpl := range folder {
		names = append(names, tmpl.Name())
	}
	return names, nil
}
func Load(directory string) (map[string]interface{}, error) {
	templateMap := make(map[string]interface{}, 0)
	names, err := List(directory)
	if err != nil {
		return nil, err
	}
	for _, tName := range names {
		content, err1 := ioutil.ReadFile(directory + string(os.PathSeparator) + tName)
		if err1 != nil {
			return nil, err1
		}
		templateMap[tName] = string(content)
	}
	return templateMap, err
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
