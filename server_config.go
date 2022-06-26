package core

type ServerConf struct {
	Name              string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Version           string `mapstructure:"version" json:"version,omitempty" gorm:"column:version" bson:"version,omitempty" dynamodbav:"version,omitempty" firestore:"version,omitempty"`
	Port              *int64 `mapstructure:"port" json:"port,omitempty" gorm:"column:port" bson:"port,omitempty" dynamodbav:"port,omitempty" firestore:"port,omitempty"`
	Secure            bool   `mapstructure:"secure" json:"secure,omitempty" gorm:"column:secure" bson:"secure,omitempty" dynamodbav:"secure,omitempty" firestore:"secure,omitempty"`
	Log               *bool  `mapstructure:"log" json:"log,omitempty" gorm:"column:log" bson:"log,omitempty" dynamodbav:"log,omitempty" firestore:"log,omitempty"`
	Monitor           *bool  `mapstructure:"monitor" json:"monitor,omitempty" gorm:"column:monitor" bson:"monitor,omitempty" dynamodbav:"monitor,omitempty" firestore:"monitor,omitempty"`
	CORS              *bool  `mapstructure:"cors" json:"cors,omitempty" gorm:"column:cors" bson:"cors,omitempty" dynamodbav:"cors,omitempty" firestore:"cors,omitempty"`
	WriteTimeout      *int64 `mapstructure:"write_timeout" json:"writeTimeout,omitempty" gorm:"column:writetimeout" bson:"writeTimeout,omitempty" dynamodbav:"writeTimeout,omitempty" firestore:"writeTimeout,omitempty"`
	ReadTimeout       *int64 `mapstructure:"read_timeout" json:"readTimeout,omitempty" gorm:"column:readtimeout" bson:"readTimeout,omitempty" dynamodbav:"readTimeout,omitempty" firestore:"readTimeout,omitempty"`
	ReadHeaderTimeout *int64 `mapstructure:"read_header_timeout" json:"readHeaderTimeout,omitempty" gorm:"column:readheadertimeout" bson:"readHeaderTimeout,omitempty" dynamodbav:"readHeaderTimeout,omitempty" firestore:"readHeaderTimeout,omitempty"`
	IdleTimeout       *int64 `mapstructure:"idle_timeout" json:"idleTimeout,omitempty" gorm:"column:idletimeout" bson:"idleTimeout,omitempty" dynamodbav:"idleTimeout,omitempty" firestore:"idleTimeout,omitempty"`
	MaxHeaderBytes    *int   `mapstructure:"max_header_bytes" json:"maxHeaderBytes,omitempty" gorm:"column:maxheaderbytes" bson:"maxHeaderBytes,omitempty" dynamodbav:"maxHeaderBytes,omitempty" firestore:"maxHeaderBytes,omitempty"`
	Cert              string `mapstructure:"cert" json:"cert,omitempty" gorm:"column:cert" bson:"cert,omitempty" dynamodbav:"cert,omitempty" firestore:"cert,omitempty"`
	Key               string `mapstructure:"key" json:"key,omitempty" gorm:"column:key" bson:"key,omitempty" dynamodbav:"key,omitempty" firestore:"key,omitempty"`
}
