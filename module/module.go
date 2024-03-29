package module

import "context"

type Module struct {
	Id    string `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"_"`
	Path  string `yaml:"path" mapstructure:"path" json:"path,omitempty" gorm:"column:path" bson:"path,omitempty" dynamodbav:"path,omitempty" firestore:"path,omitempty"`
	Route string `yaml:"route" mapstructure:"route" json:"route,omitempty" gorm:"column:route" bson:"route,omitempty" dynamodbav:"route,omitempty" firestore:"route,omitempty"`
}

type Loader interface {
	Load(ctx context.Context) ([]Module, error)
}
