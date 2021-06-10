package diff

type DiffModel struct {
	Id     interface{} `mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Origin interface{} `mapstructure:"origin" json:"origin,omitempty" gorm:"column:origin" bson:"origin,omitempty" dynamodbav:"origin,omitempty" firestore:"origin,omitempty"`
	Value  interface{} `mapstructure:"value" json:"value,omitempty" gorm:"column:value" bson:"value,omitempty" dynamodbav:"value,omitempty" firestore:"value,omitempty"`
	By     string      `mapstructure:"by" json:"by,omitempty" gorm:"column:updated_by" bson:"by,omitempty" dynamodbav:"by,omitempty" firestore:"by,omitempty"`
}
