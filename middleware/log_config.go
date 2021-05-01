package middleware

type LogConfig struct {
	Separate       bool              `mapstructure:"separate" json:"separate,omitempty" gorm:"column:separate" bson:"separate,omitempty" dynamodbav:"separate,omitempty" firestore:"separate,omitempty"`
	Build          bool              `mapstructure:"build" json:"build,omitempty" gorm:"column:build" bson:"build,omitempty" dynamodbav:"build,omitempty" firestore:"build,omitempty"`
	Log            bool              `mapstructure:"log" json:"log,omitempty" gorm:"column:log" bson:"log,omitempty" dynamodbav:"log,omitempty" firestore:"log,omitempty"`
	Skips          string            `mapstructure:"skips" json:"skips,omitempty" gorm:"column:skips" bson:"skips,omitempty" dynamodbav:"skips,omitempty" firestore:"skips,omitempty"`
	Ip             string            `mapstructure:"ip" json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	Duration       string            `mapstructure:"duration" json:"duration,omitempty" gorm:"column:duration" bson:"duration,omitempty" dynamodbav:"duration,omitempty" firestore:"duration,omitempty"`
	Uri            string            `mapstructure:"uri" json:"uri,omitempty" gorm:"column:uri" bson:"uri,omitempty" dynamodbav:"uri,omitempty" firestore:"uri,omitempty"`
	Body           string            `mapstructure:"body" json:"body,omitempty" gorm:"column:body" bson:"body,omitempty" dynamodbav:"body,omitempty" firestore:"body,omitempty"`
	Size           string            `mapstructure:"size" json:"size,omitempty" gorm:"column:size" bson:"size,omitempty" dynamodbav:"size,omitempty" firestore:"size,omitempty"`
	ReqId          string            `mapstructure:"req_id" json:"reqId,omitempty" gorm:"column:reqid" bson:"reqId,omitempty" dynamodbav:"reqId,omitempty" firestore:"reqId,omitempty"`
	Scheme         string            `mapstructure:"scheme" json:"scheme,omitempty" gorm:"column:scheme" bson:"scheme,omitempty" dynamodbav:"scheme,omitempty" firestore:"scheme,omitempty"`
	Proto          string            `mapstructure:"proto" json:"proto,omitempty" gorm:"column:proto" bson:"proto,omitempty" dynamodbav:"proto,omitempty" firestore:"proto,omitempty"`
	Method         string            `mapstructure:"method" json:"method,omitempty" gorm:"column:method" bson:"method,omitempty" dynamodbav:"method,omitempty" firestore:"method,omitempty"`
	RemoteAddr     string            `mapstructure:"remote_addr" json:"remoteAddr,omitempty" gorm:"column:remoteAddr" bson:"remoteAddr,omitempty" dynamodbav:"remoteAddr,omitempty" firestore:"remoteAddr,omitempty"`
	RemoteIp       string            `mapstructure:"remote_ip" json:"remoteIp,omitempty" gorm:"column:remoteIp" bson:"remoteIp,omitempty" dynamodbav:"remoteIp,omitempty" firestore:"remoteIp,omitempty"`
	UserAgent      string            `mapstructure:"user_agent" json:"userAgent,omitempty" gorm:"column:userAgent" bson:"userAgent,omitempty" dynamodbav:"userAgent,omitempty" firestore:"userAgent,omitempty"`
	ResponseStatus string            `mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	Request        string            `mapstructure:"request" json:"request,omitempty" gorm:"column:request" bson:"request,omitempty" dynamodbav:"request,omitempty" firestore:"request,omitempty"`
	Response       string            `mapstructure:"response" json:"response,omitempty" gorm:"column:response" bson:"response,omitempty" dynamodbav:"response,omitempty" firestore:"response,omitempty"`
	Fields         string            `mapstructure:"fields" json:"fields,omitempty" gorm:"column:fields" bson:"fields,omitempty" dynamodbav:"fields,omitempty" firestore:"fields,omitempty"`
	Masks          string            `mapstructure:"masks" json:"masks,omitempty" gorm:"column:masks" bson:"masks,omitempty" dynamodbav:"masks,omitempty" firestore:"masks,omitempty"`
	Map            map[string]string `mapstructure:"map" json:"map,omitempty" gorm:"column:map" bson:"map,omitempty" dynamodbav:"map,omitempty" firestore:"map,omitempty"`
	Constants      map[string]string `mapstructure:"constants" json:"constants,omitempty" gorm:"column:constants" bson:"constants,omitempty" dynamodbav:"constants,omitempty" firestore:"constants,omitempty"`
}

type FieldConfig struct {
	Log       bool              `mapstructure:"log" json:"log,omitempty" gorm:"column:log" bson:"log,omitempty" dynamodbav:"log,omitempty" firestore:"log,omitempty"`
	Ip        string            `mapstructure:"ip" json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	Map       map[string]string `mapstructure:"map" json:"map,omitempty" gorm:"column:map" bson:"map,omitempty" dynamodbav:"map,omitempty" firestore:"map,omitempty"`
	Constants map[string]string `mapstructure:"constants" json:"constants,omitempty" gorm:"column:constants" bson:"constants,omitempty" dynamodbav:"constants,omitempty" firestore:"constants,omitempty"`
	Duration  string            `mapstructure:"duration" json:"duration,omitempty" gorm:"column:duration" bson:"duration,omitempty" dynamodbav:"duration,omitempty" firestore:"duration,omitempty"`
	Fields    []string          `mapstructure:"fields" json:"fields,omitempty" gorm:"column:fields" bson:"fields,omitempty" dynamodbav:"fields,omitempty" firestore:"fields,omitempty"`
	Masks     []string          `mapstructure:"masks" json:"masks,omitempty" gorm:"column:masks" bson:"masks,omitempty" dynamodbav:"masks,omitempty" firestore:"masks,omitempty"`
	Skips     []string          `mapstructure:"skips" json:"skips,omitempty" gorm:"column:skips" bson:"skips,omitempty" dynamodbav:"skips,omitempty" firestore:"skips,omitempty"`
}
