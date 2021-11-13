package cors

import (
	"github.com/rs/cors"
	"strings"
)

type AllowConfig struct {
	Origins            string `mapstructure:"origins" json:"origins,omitempty" gorm:"column:origins" bson:"origins,omitempty" dynamodbav:"origins,omitempty" firestore:"origins,omitempty"`
	Methods            string `mapstructure:"methods" json:"methods,omitempty" gorm:"column:methods" bson:"methods,omitempty" dynamodbav:"methods,omitempty" firestore:"methods,omitempty"`
	Headers            string `mapstructure:"headers" json:"headers,omitempty" gorm:"column:headers" bson:"headers,omitempty" dynamodbav:"headers,omitempty" firestore:"headers,omitempty"`
	Credentials        bool   `mapstructure:"credentials" json:"credentials,omitempty" gorm:"column:credentials" bson:"credentials,omitempty" dynamodbav:"credentials,omitempty" firestore:"credentials,omitempty"`
	MaxAge             *int   `mapstructure:"max_age" json:"maxAge,omitempty" gorm:"column:maxage" bson:"maxAge,omitempty" dynamodbav:"maxAge,omitempty" firestore:"maxAge,omitempty"`
	ExposedHeaders     string `mapstructure:"exposed_headers" json:"exposedHeaders,omitempty" gorm:"column:exposedheaders" bson:"exposedHeaders,omitempty" dynamodbav:"exposedHeaders,omitempty" firestore:"exposedHeaders,omitempty"`
	OptionsPassthrough *bool  `mapstructure:"options_passthrough" json:"optionsPassthrough,omitempty" gorm:"column:optionsPassthrough" bson:"optionsPassthrough,omitempty" dynamodbav:"optionsPassthrough,omitempty" firestore:"optionsPassthrough,omitempty"`
}

func New(conf AllowConfig) *cors.Cors {
	opts := cors.Options{AllowCredentials: conf.Credentials}
	if len(conf.Headers) > 0 {
		opts.AllowedHeaders = strings.Split(conf.Headers, ",")
	}
	if len(conf.Origins) > 0 {
		opts.AllowedOrigins = strings.Split(conf.Origins, ",")
	}
	if len(conf.Methods) > 0 {
		opts.AllowedMethods = strings.Split(conf.Methods, ",")
	}
	if conf.MaxAge != nil {
		opts.MaxAge = *conf.MaxAge
	}
	if len(conf.ExposedHeaders) > 0 {
		opts.ExposedHeaders = strings.Split(conf.ExposedHeaders, ",")
	}
	if conf.OptionsPassthrough != nil {
		opts.OptionsPassthrough = *conf.OptionsPassthrough
	}
	return cors.New(opts)
}
