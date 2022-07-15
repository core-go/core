package echo

type LogConfig struct {
	Separate       bool              `yaml:"separate" mapstructure:"separate" json:"separate,omitempty" gorm:"column:separate" bson:"separate,omitempty" dynamodbav:"separate,omitempty" firestore:"separate,omitempty"`
	Build          bool              `yaml:"build" mapstructure:"build" json:"build,omitempty" gorm:"column:build" bson:"build,omitempty" dynamodbav:"build,omitempty" firestore:"build,omitempty"`
	Log            bool              `yaml:"log" mapstructure:"log" json:"log,omitempty" gorm:"column:log" bson:"log,omitempty" dynamodbav:"log,omitempty" firestore:"log,omitempty"`
	Skips          string            `yaml:"skips" mapstructure:"skips" json:"skips,omitempty" gorm:"column:skips" bson:"skips,omitempty" dynamodbav:"skips,omitempty" firestore:"skips,omitempty"`
	Ip             string            `yaml:"ip" mapstructure:"ip" json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	Duration       string            `yaml:"duration" mapstructure:"duration" json:"duration,omitempty" gorm:"column:duration" bson:"duration,omitempty" dynamodbav:"duration,omitempty" firestore:"duration,omitempty"`
	Uri            string            `yaml:"uri" mapstructure:"uri" json:"uri,omitempty" gorm:"column:uri" bson:"uri,omitempty" dynamodbav:"uri,omitempty" firestore:"uri,omitempty"`
	Body           string            `yaml:"body" mapstructure:"body" json:"body,omitempty" gorm:"column:body" bson:"body,omitempty" dynamodbav:"body,omitempty" firestore:"body,omitempty"`
	Size           string            `yaml:"size" mapstructure:"size" json:"size,omitempty" gorm:"column:size" bson:"size,omitempty" dynamodbav:"size,omitempty" firestore:"size,omitempty"`
	ReqId          string            `yaml:"req_id" mapstructure:"req_id" json:"reqId,omitempty" gorm:"column:reqid" bson:"reqId,omitempty" dynamodbav:"reqId,omitempty" firestore:"reqId,omitempty"`
	Scheme         string            `yaml:"scheme" mapstructure:"scheme" json:"scheme,omitempty" gorm:"column:scheme" bson:"scheme,omitempty" dynamodbav:"scheme,omitempty" firestore:"scheme,omitempty"`
	Proto          string            `yaml:"proto" mapstructure:"proto" json:"proto,omitempty" gorm:"column:proto" bson:"proto,omitempty" dynamodbav:"proto,omitempty" firestore:"proto,omitempty"`
	Method         string            `yaml:"method" mapstructure:"method" json:"method,omitempty" gorm:"column:method" bson:"method,omitempty" dynamodbav:"method,omitempty" firestore:"method,omitempty"`
	RemoteAddr     string            `yaml:"remote_addr" mapstructure:"remote_addr" json:"remoteAddr,omitempty" gorm:"column:remoteAddr" bson:"remoteAddr,omitempty" dynamodbav:"remoteAddr,omitempty" firestore:"remoteAddr,omitempty"`
	RemoteIp       string            `yaml:"remote_ip" mapstructure:"remote_ip" json:"remoteIp,omitempty" gorm:"column:remoteIp" bson:"remoteIp,omitempty" dynamodbav:"remoteIp,omitempty" firestore:"remoteIp,omitempty"`
	UserAgent      string            `yaml:"user_agent" mapstructure:"user_agent" json:"userAgent,omitempty" gorm:"column:userAgent" bson:"userAgent,omitempty" dynamodbav:"userAgent,omitempty" firestore:"userAgent,omitempty"`
	ResponseStatus string            `yaml:"status" mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	Request        string            `yaml:"request" mapstructure:"request" json:"request,omitempty" gorm:"column:request" bson:"request,omitempty" dynamodbav:"request,omitempty" firestore:"request,omitempty"`
	Response       string            `yaml:"response" mapstructure:"response" json:"response,omitempty" gorm:"column:response" bson:"response,omitempty" dynamodbav:"response,omitempty" firestore:"response,omitempty"`
	Fields         string            `yaml:"fields" mapstructure:"fields" json:"fields,omitempty" gorm:"column:fields" bson:"fields,omitempty" dynamodbav:"fields,omitempty" firestore:"fields,omitempty"`
	Masks          string            `yaml:"masks" mapstructure:"masks" json:"masks,omitempty" gorm:"column:masks" bson:"masks,omitempty" dynamodbav:"masks,omitempty" firestore:"masks,omitempty"`
	Map            map[string]string `yaml:"map" mapstructure:"map" json:"map,omitempty" gorm:"column:map" bson:"map,omitempty" dynamodbav:"map,omitempty" firestore:"map,omitempty"`
	Constants      map[string]string `yaml:"constants" mapstructure:"constants" json:"constants,omitempty" gorm:"column:constants" bson:"constants,omitempty" dynamodbav:"constants,omitempty" firestore:"constants,omitempty"`
}

type FieldConfig struct {
	Log       bool              `yaml:"log" mapstructure:"log" json:"log,omitempty" gorm:"column:log" bson:"log,omitempty" dynamodbav:"log,omitempty" firestore:"log,omitempty"`
	Ip        string            `yaml:"ip" mapstructure:"ip" json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	Map       map[string]string `yaml:"map" mapstructure:"map" json:"map,omitempty" gorm:"column:map" bson:"map,omitempty" dynamodbav:"map,omitempty" firestore:"map,omitempty"`
	Constants map[string]string `yaml:"constants" mapstructure:"constants" json:"constants,omitempty" gorm:"column:constants" bson:"constants,omitempty" dynamodbav:"constants,omitempty" firestore:"constants,omitempty"`
	Duration  string            `yaml:"duration" mapstructure:"duration" json:"duration,omitempty" gorm:"column:duration" bson:"duration,omitempty" dynamodbav:"duration,omitempty" firestore:"duration,omitempty"`
	Fields    []string          `yaml:"fields" mapstructure:"fields" json:"fields,omitempty" gorm:"column:fields" bson:"fields,omitempty" dynamodbav:"fields,omitempty" firestore:"fields,omitempty"`
	Masks     []string          `yaml:"masks" mapstructure:"masks" json:"masks,omitempty" gorm:"column:masks" bson:"masks,omitempty" dynamodbav:"masks,omitempty" firestore:"masks,omitempty"`
	Skips     []string          `yaml:"skips" mapstructure:"skips" json:"skips,omitempty" gorm:"column:skips" bson:"skips,omitempty" dynamodbav:"skips,omitempty" firestore:"skips,omitempty"`
}
