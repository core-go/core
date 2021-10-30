package io

import "os"

type File struct {
	Name    string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Content string `mapstructure:"content" json:"content,omitempty" gorm:"column:content" bson:"content,omitempty" dynamodbav:"content,omitempty" firestore:"content,omitempty"`
}

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
