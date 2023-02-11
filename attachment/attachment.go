package attachment

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Attachment struct {
	OriginalFileName string `json:"originalFileName,omitempty" gorm:"column:original_filename" bson:"original_filename,omitempty" dynamodbav:"original_filename,omitempty" firestore:"original_filename,omitempty"`
	FileName         string `json:"fileName,omitempty" gorm:"column:fileName" bson:"fileName,omitempty" dynamodbav:"fileName,omitempty" firestore:"fileName,omitempty"`
	Url              string `json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty" avro:"url"`
	Type             string `json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
	Size             int64  `json:"size,omitempty" gorm:"column:type" bson:"size,omitempty" dynamodbav:"size,omitempty" firestore:"size,omitempty"`
}

func (h Attachment) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *Attachment) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &h)
}
