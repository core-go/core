package service

type DiffModel struct {
	Id       string `json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	OldValue string `json:"oldValue,omitempty" gorm:"column:oldvalue" bson:"oldValue,omitempty" dynamodbav:"oldValue,omitempty" firestore:"oldValue,omitempty"`
	NewValue string `json:"newValue,omitempty" gorm:"column:newvalue" bson:"newValue,omitempty" dynamodbav:"newValue,omitempty" firestore:"newValue,omitempty"`
}
