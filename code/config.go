package code

type Config struct {
	Handler HandlerConfig   `mapstructure:"handler" json:"handler,omitempty" gorm:"column:handler" bson:"handler,omitempty" dynamodbav:"handler,omitempty" firestore:"handler,omitempty"`
	Loader  StructureConfig `mapstructure:"loader" json:"loader,omitempty" gorm:"column:loader" bson:"loader,omitempty" dynamodbav:"loader,omitempty" firestore:"loader,omitempty"`
}
