package ldap

type LdapConfig struct {
	Server             string            `yaml:"server" mapstructure:"server" json:"server,omitempty" gorm:"column:server" bson:"server,omitempty" dynamodbav:"server,omitempty" firestore:"server,omitempty"`
	BaseDN             string            `yaml:"base_dn" mapstructure:"base_dn" json:"baseDN,omitempty" gorm:"column:basedn" bson:"baseDN,omitempty" dynamodbav:"baseDN,omitempty" firestore:"baseDN,omitempty"`
	Timeout            int64             `yaml:"timeout" mapstructure:"timeout" json:"timeout,omitempty" gorm:"column:timeout" bson:"timeout,omitempty" dynamodbav:"timeout,omitempty" firestore:"timeout,omitempty"`
	Domain             string            `yaml:"domain" mapstructure:"domain" json:"domain,omitempty" gorm:"column:domain" bson:"domain,omitempty" dynamodbav:"domain,omitempty" firestore:"domain,omitempty"`
	Username           string            `yaml:"username" mapstructure:"username" json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Password           string            `yaml:"password" mapstructure:"password" json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
	Filter             string            `yaml:"filter" mapstructure:"filter" json:"filter,omitempty" gorm:"column:filter" bson:"filter,omitempty" dynamodbav:"filter,omitempty" firestore:"filter,omitempty"`
	TLS                bool              `yaml:"tls" mapstructure:"tls" json:"tls,omitempty" gorm:"column:tls" bson:"tls,omitempty" dynamodbav:"tls,omitempty" firestore:"tls,omitempty"`
	StartTLS           bool              `yaml:"start_tls" mapstructure:"start_tls" json:"startTLS,omitempty" gorm:"column:starttls" bson:"startTLS,omitempty" dynamodbav:"startTLS,omitempty" firestore:"startTLS,omitempty"`
	InsecureSkipVerify bool              `yaml:"insecure_skip_verify" mapstructure:"insecure_skip_verify" json:"insecureSkipVerify,omitempty" gorm:"column:insecureskipverify" bson:"insecureSkipVerify,omitempty" dynamodbav:"insecureSkipVerify,omitempty" firestore:"insecureSkipVerify,omitempty"`
	Attributes         map[string]string `yaml:"attributes" mapstructure:"attributes" json:"attributes,omitempty" gorm:"column:attributes" bson:"attributes,omitempty" dynamodbav:"attributes,omitempty" firestore:"attributes,omitempty"`
	Dates              map[string]string `yaml:"dates" mapstructure:"dates" json:"dates,omitempty" gorm:"column:dates" bson:"dates,omitempty" dynamodbav:"dates,omitempty" firestore:"dates,omitempty"`
}
