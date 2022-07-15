package mail

type SimpleMail struct {
	From    Email     `yaml:"from" mapstructure:"from" json:"from,omitempty" gorm:"column:from" bson:"from,omitempty" dynamodbav:"from,omitempty" firestore:"from,omitempty"`
	To      []Email   `yaml:"to" mapstructure:"to" json:"to,omitempty" gorm:"column:to" bson:"to,omitempty" dynamodbav:"to,omitempty" firestore:"to,omitempty"`
	Cc      *[]Email  `yaml:"cc" mapstructure:"cc" json:"cc,omitempty" gorm:"column:cc" bson:"cc,omitempty" dynamodbav:"cc,omitempty" firestore:"cc,omitempty"`
	Bcc     *[]Email  `yaml:"bcc" mapstructure:"bcc" json:"bcc,omitempty" gorm:"column:bcc" bson:"bcc,omitempty" dynamodbav:"bcc,omitempty" firestore:"bcc,omitempty"`
	Subject string    `yaml:"subject" mapstructure:"subject" json:"subject,omitempty" gorm:"column:subject" bson:"subject,omitempty" dynamodbav:"subject,omitempty" firestore:"subject,omitempty"`
	Content []Content `yaml:"content" mapstructure:"content" json:"content,omitempty" gorm:"column:content" bson:"content,omitempty" dynamodbav:"content,omitempty" firestore:"content,omitempty"`
}
