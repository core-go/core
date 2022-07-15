package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ActivityLogSchemaConfig struct {
	User      string    `yaml:"user" mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	Ip        string    `yaml:"ip" mapstructure:"ip" json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	Resource  string    `yaml:"resource" mapstructure:"resource" json:"resource,omitempty" gorm:"column:resource" bson:"resource,omitempty" dynamodbav:"resource,omitempty" firestore:"resource,omitempty"`
	Action    string    `yaml:"action" mapstructure:"action" json:"action,omitempty" gorm:"column:action" bson:"action,omitempty" dynamodbav:"action,omitempty" firestore:"action,omitempty"`
	Timestamp string    `yaml:"timestamp" mapstructure:"timestamp" json:"timestamp,omitempty" gorm:"column:timestamp" bson:"timestamp,omitempty" dynamodbav:"timestamp,omitempty" firestore:"timestamp,omitempty"`
	Status    string    `yaml:"status" mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	Desc      string    `yaml:"desc" mapstructure:"desc" json:"desc,omitempty" gorm:"column:desc" bson:"desc,omitempty" dynamodbav:"desc,omitempty" firestore:"desc,omitempty"`
	Ext       *[]string `yaml:"ext" mapstructure:"ext" json:"ext,omitempty" gorm:"column:ext" bson:"ext,omitempty" dynamodbav:"ext,omitempty" firestore:"ext,omitempty"`
}

func NewActivityLogWriter(database *mongo.Database, collectionName string, config ActivityLogConfig, schema ActivityLogSchemaConfig, generate func(context.Context) (string, error)) *ActivityLogWriter {
	if len(schema.User) == 0 {
		schema.User = "user"
	}
	if len(schema.Resource) == 0 {
		schema.Resource = "resource"
	}
	if len(schema.Action) == 0 {
		schema.Action = "action"
	}
	if len(schema.Timestamp) == 0 {
		schema.Timestamp = "timestamp"
	}
	if len(schema.Status) == 0 {
		schema.Status = "status"
	}
	if len(schema.Desc) == 0 {
		schema.Desc = "desc"
	}
	col := database.Collection(collectionName)
	sender := ActivityLogWriter{Database: database, Collection: col, Config: config, Schema: schema, Generate: generate}
	return &sender
}

type ActivityLogWriter struct {
	Database   *mongo.Database
	Collection *mongo.Collection
	Config     ActivityLogConfig
	Schema     ActivityLogSchemaConfig
	Generate   func(ctx context.Context) (string, error)
}

type ActivityLogConfig struct {
	User       string `yaml:"user" mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	Ip         string `yaml:"ip" mapstructure:"ip" json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	True       string `yaml:"true" mapstructure:"true" json:"true,omitempty" gorm:"column:true" bson:"true,omitempty" dynamodbav:"true,omitempty" firestore:"true,omitempty"`
	False      string `yaml:"false" mapstructure:"false" json:"false,omitempty" gorm:"column:false" bson:"false,omitempty" dynamodbav:"false,omitempty" firestore:"false,omitempty"`
	Goroutines bool   `yaml:"goroutines" mapstructure:"goroutines" json:"goroutines,omitempty" gorm:"column:goroutines" bson:"goroutines,omitempty" dynamodbav:"goroutines,omitempty" firestore:"goroutines,omitempty"`
}

func (s *ActivityLogWriter) Write(ctx context.Context, resource string, action string, success bool, desc string) error {
	log := make(map[string]interface{})
	now := time.Now()
	ch := s.Schema
	log[ch.Timestamp] = &now
	log[ch.Resource] = resource
	log[ch.Action] = action
	log[ch.Desc] = desc

	if success {
		log[ch.Status] = s.Config.True
	} else {
		log[ch.Status] = s.Config.False
	}
	log[ch.User] = getString(ctx, s.Config.User)
	if len(ch.Ip) > 0 {
		log[ch.Ip] = getString(ctx, s.Config.Ip)
	}
	if s.Generate != nil {
		id, er0 := s.Generate(ctx)
		if er0 == nil && len(id) > 0 {
			log["_id"] = id
		}
	}
	ext := buildExt(ctx, ch.Ext)
	if len(ext) > 0 {
		for k, v := range ext {
			log[k] = v
		}
	}
	if !s.Config.Goroutines {
		_, er3 := s.Collection.InsertOne(ctx, log)
		return er3
	} else {
		go s.Collection.InsertOne(ctx, log)
		return nil
	}
}

func buildExt(ctx context.Context, keys *[]string) map[string]interface{} {
	headers := make(map[string]interface{})
	if keys != nil {
		hs := *keys
		for _, header := range hs {
			v := ctx.Value(header)
			if v != nil {
				headers[header] = v
			}
		}
	}
	return headers
}
func getString(ctx context.Context, key string) string {
	if len(key) > 0 {
		u := ctx.Value(key)
		if u != nil {
			s, ok := u.(string)
			if ok {
				return s
			} else {
				return ""
			}
		}
	}
	return ""
}
