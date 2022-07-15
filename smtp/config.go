package smtp

import "github.com/core-go/core/mail"

type MailConfig struct {
	Provider string       `yaml:"provider" mapstructure:"provider" json:"provider,omitempty" gorm:"column:provider" bson:"provider,omitempty" dynamodbav:"provider,omitempty" firestore:"provider,omitempty"`
	From     mail.Email   `yaml:"from" mapstructure:"from" json:"from,omitempty" gorm:"column:from" bson:"from,omitempty" dynamodbav:"from,omitempty" firestore:"from,omitempty"`
	SMTP     DialerConfig `yaml:"smtp" mapstructure:"smtp" json:"smtp,omitempty" gorm:"column:smtp" bson:"smtp,omitempty" dynamodbav:"smtp,omitempty" firestore:"smtp,omitempty"`
	APIkey   string       `yaml:"api_key" mapstructure:"api_key" json:"apiKey,omitempty" gorm:"column:apikey" bson:"apiKey,omitempty" dynamodbav:"apiKey,omitempty" firestore:"apiKey,omitempty"`
}
