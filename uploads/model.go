package uploads

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Request struct {
	OriginalFileName string `yaml:"original_file_name" mapstructure:"original_file_name" json:"originalFileName,omitempty" gorm:"column:original_file_name" bson:"original_filename,omitempty" dynamodbav:"original_filename,omitempty" firestore:"original_filename,omitempty"`
	Filename         string `yaml:"name" mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Type             string `yaml:"type" mapstructure:"type" json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
	Size             int64  `yaml:"size" mapstructure:"size" json:"size,omitempty" gorm:"column:size" bson:"size,omitempty" dynamodbav:"size,omitempty" firestore:"size,omitempty"`
	Data             []byte `yaml:"data" mapstructure:"data" json:"data,omitempty" gorm:"column:data" bson:"data,omitempty" dynamodbav:"data,omitempty" firestore:"data,omitempty"`
}

type Upload struct {
	OriginalFileName string `yaml:"original_file_name" mapstructure:"original_file_name" json:"originalFileName,omitempty" gorm:"column:original_file_name" bson:"original_filename,omitempty" dynamodbav:"original_filename,omitempty" firestore:"original_filename,omitempty"`
	FileName         string `yaml:"file_name" mapstructure:"file_name" json:"fileName,omitempty" gorm:"column:file_name" bson:"fileName,omitempty" dynamodbav:"fileName,omitempty" firestore:"fileName,omitempty"`
	Url              string `yaml:"url" mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty" avro:"url" validate:"required"`
	Type             string `yaml:"type" mapstructure:"type" json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
	Size             int64  `yaml:"size" mapstructure:"size" json:"size,omitempty" gorm:"column:size" bson:"size,omitempty" dynamodbav:"size,omitempty" firestore:"size,omitempty"`
}

func (u Upload) Value() (driver.Value, error) {
	return json.Marshal(u)
}
func (u *Upload) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &u)
}

type FileConfig struct {
	MaxSize           int64  `yaml:"max_size" mapstructure:"max_size" json:"maxSize,omitempty"`
	MaxSizeMemory     int64  `yaml:"max_size_memory" mapstructure:"max_size_memory" json:"maxSizeMemory,omitempty"`
	AllowedExtensions string `yaml:"allowed_extensions" mapstructure:"allowed_extensions" json:"allowedExtensions,omitempty"`
}
