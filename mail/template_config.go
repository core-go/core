package mail

type TemplateConfig struct {
	Subject string `mapstructure:"subject" json:"subject,omitempty" gorm:"column:subject" bson:"subject,omitempty" dynamodbav:"subject,omitempty" firestore:"subject,omitempty"`
	Body    string `mapstructure:"body" json:"body,omitempty" gorm:"column:body" bson:"body,omitempty" dynamodbav:"body,omitempty" firestore:"body,omitempty"`
}
