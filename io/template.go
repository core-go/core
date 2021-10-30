package io

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	Name    string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Content string `mapstructure:"content" json:"content,omitempty" gorm:"column:content" bson:"content,omitempty" dynamodbav:"content,omitempty" firestore:"content,omitempty"`
}

func ListFileNames(path string) ([]string, error) {
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

func LoadTemplates(templateFolder string) (map[string]interface{}, error) {
	templateMap := make(map[string]interface{}, 0)
	names, err := ListFileNames(templateFolder)
	if err != nil {
		return nil, err
	}
	for _, tName := range names {
		content, err1 := ioutil.ReadFile(templateFolder + string(os.PathSeparator) + tName)
		if err1 != nil {
			return nil, err1
		}
		templateMap[tName] = string(content)
	}
	return templateMap, err
}
