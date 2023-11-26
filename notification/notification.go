package notification

import "context"

type Notification struct {
	Id       string `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Sender   string `yaml:"sender" mapstructure:"sender" json:"sender,omitempty" gorm:"column:sender" bson:"sender,omitempty" dynamodbav:"sender,omitempty" firestore:"sender,omitempty"`
	Receiver string `yaml:"receiver" mapstructure:"receiver" json:"receiver,omitempty" gorm:"column:receiver" bson:"receiver,omitempty" dynamodbav:"receiver,omitempty" firestore:"receiver,omitempty"`
	Url      string `yaml:"url" mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
	Message  string `yaml:"message" mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}

func Build(sender string, receiver string, url string, message string) *Notification {
	return &Notification{
		Sender: sender,
		Receiver: receiver,
		Url: url,
		Message: message,
	}
}
type NotificationPort interface {
	Push(ctx context.Context, noti *Notification) (int64, error)
	PushNotifications(ctx context.Context, notifications []Notification) (int64, error)
}
