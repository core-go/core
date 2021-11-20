package code

type Config struct {
	Handler HandlerConfig   `yaml:"handler" mapstructure:"handler" json:"handler,omitempty" gorm:"column:handler" bson:"handler,omitempty" dynamodbav:"handler,omitempty" firestore:"handler,omitempty"`
	Loader  StructureConfig `yaml:"loader" mapstructure:"loader" json:"loader,omitempty" gorm:"column:loader" bson:"loader,omitempty" dynamodbav:"loader,omitempty" firestore:"loader,omitempty"`
}
