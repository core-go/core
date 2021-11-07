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
func Load(directory string) (map[string]string, error) {
	tm := make(map[string]string, 0)
	names, er1 := List(directory)
	if er1 != nil {
		return nil, er1
	}
	for _, name := range names {
		content, er2 := ioutil.ReadFile(directory + string(os.PathSeparator) + name)
		if er2 != nil {
			return nil, er2
		}
		tm[name] = string(content)
	}
	return tm, nil
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
