package service

import v "github.com/common-go/validator"

type ResultInfo struct {
	Status  int              `mapstructure:"status" json:"status" gorm:"column:status" bson:"status" dynamodbav:"status" firestore:"status"`
	Errors  []v.ErrorMessage `mapstructure:"errors" json:"errors,omitempty" gorm:"column:errors" bson:"errors,omitempty" dynamodbav:"errors,omitempty" firestore:"errors,omitempty"`
	Value   interface{}      `mapstructure:"value" json:"value,omitempty" gorm:"column:value" bson:"value,omitempty" dynamodbav:"value,omitempty" firestore:"value,omitempty"`
	Message string           `mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}
