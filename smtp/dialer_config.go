package smtp

type DialerConfig struct {
	Host     string `yaml:"host" mapstructure:"host" json:"host,omitempty" gorm:"column:host" bson:"host,omitempty" dynamodbav:"host,omitempty" firestore:"host,omitempty"`
	Port     int    `yaml:"port" mapstructure:"port" json:"port,omitempty" gorm:"column:port" bson:"port,omitempty" dynamodbav:"port,omitempty" firestore:"port,omitempty"`
	Username string `yaml:"username" mapstructure:"username" json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Password string `yaml:"password" mapstructure:"password" json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
	SSL      *bool  `yaml:"ssl" mapstructure:"ssl" json:"ssl,omitempty" gorm:"column:ssl" bson:"ssl,omitempty" dynamodbav:"ssl,omitempty" firestore:"ssl,omitempty"`
}
