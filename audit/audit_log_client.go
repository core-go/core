package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type ClientConfig struct {
	Url    string         `yaml:"url" mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
	Schema AuditLogSchema `yaml:"schema" mapstructure:"schema" json:"schema,omitempty" gorm:"column:schema" bson:"schema,omitempty" dynamodbav:"schema,omitempty" firestore:"schema,omitempty"`
	Config AuditLogConfig `yaml:"config" mapstructure:"config" json:"config,omitempty" gorm:"column:config" bson:"config,omitempty" dynamodbav:"config,omitempty" firestore:"config,omitempty"`
	Retry  Retry          `yaml:"retry" mapstructure:"retry" json:"retry,omitempty" gorm:"column:retry" bson:"retry,omitempty" dynamodbav:"retry,omitempty" firestore:"retry,omitempty"`
}
type Retry struct {
	Retry1  int64 `yaml:"1" mapstructure:"1" json:"retry1,omitempty" gorm:"column:retry1" bson:"retry1,omitempty" dynamodbav:"retry1,omitempty" firestore:"retry1,omitempty"`
	Retry2  int64 `yaml:"2" mapstructure:"2" json:"retry2,omitempty" gorm:"column:retry2" bson:"retry2,omitempty" dynamodbav:"retry2,omitempty" firestore:"retry2,omitempty"`
	Retry3  int64 `yaml:"3" mapstructure:"3" json:"retry3,omitempty" gorm:"column:retry3" bson:"retry3,omitempty" dynamodbav:"retry3,omitempty" firestore:"retry3,omitempty"`
	Retry4  int64 `yaml:"4" mapstructure:"4" json:"retry4,omitempty" gorm:"column:retry4" bson:"retry4,omitempty" dynamodbav:"retry4,omitempty" firestore:"retry4,omitempty"`
	Retry5  int64 `yaml:"5" mapstructure:"5" json:"retry5,omitempty" gorm:"column:retry5" bson:"retry5,omitempty" dynamodbav:"retry5,omitempty" firestore:"retry5,omitempty"`
	Retry6  int64 `yaml:"6" mapstructure:"6" json:"retry6,omitempty" gorm:"column:retry6" bson:"retry6,omitempty" dynamodbav:"retry6,omitempty" firestore:"retry6,omitempty"`
	Retry7  int64 `yaml:"7" mapstructure:"7" json:"retry7,omitempty" gorm:"column:retry7" bson:"retry7,omitempty" dynamodbav:"retry7,omitempty" firestore:"retry7,omitempty"`
	Retry8  int64 `yaml:"8" mapstructure:"8" json:"retry8,omitempty" gorm:"column:retry8" bson:"retry8,omitempty" dynamodbav:"retry8,omitempty" firestore:"retry8,omitempty"`
	Retry9  int64 `yaml:"9" mapstructure:"9" json:"retry9,omitempty" gorm:"column:retry9" bson:"retry9,omitempty" dynamodbav:"retry9,omitempty" firestore:"retry9,omitempty"`
	Retry10 int64 `yaml:"10" mapstructure:"10" json:"retry10,omitempty" gorm:"column:retry10" bson:"retry10,omitempty" dynamodbav:"retry10,omitempty" firestore:"retry10,omitempty"`
	Retry11 int64 `yaml:"11" mapstructure:"11" json:"retry11,omitempty" gorm:"column:retry11" bson:"retry11,omitempty" dynamodbav:"retry11,omitempty" firestore:"retry11,omitempty"`
	Retry12 int64 `yaml:"12" mapstructure:"12" json:"retry12,omitempty" gorm:"column:retry12" bson:"retry12,omitempty" dynamodbav:"retry12,omitempty" firestore:"retry12,omitempty"`
}

func NewAuditLogClient(client *http.Client, url string, config AuditLogConfig, schema AuditLogSchema, generate func(context.Context) (string, error), logError func(context.Context, string), transform func(map[string]interface{}) map[string]interface{}, retries ...time.Duration) *AuditLogClient {
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
	sender := AuditLogClient{Client: client, Url: url, Config: config, Schema: schema, Generate: generate, Transform: transform, Error: logError, Retries: retries}
	return &sender
}

type AuditLogConfig struct {
	User       string `yaml:"user" mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	Ip         string `yaml:"ip" mapstructure:"ip" json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	True       string `yaml:"true" mapstructure:"true" json:"true,omitempty" gorm:"column:true" bson:"true,omitempty" dynamodbav:"true,omitempty" firestore:"true,omitempty"`
	False      string `yaml:"false" mapstructure:"false" json:"false,omitempty" gorm:"column:false" bson:"false,omitempty" dynamodbav:"false,omitempty" firestore:"false,omitempty"`
	Goroutines bool   `yaml:"goroutines" mapstructure:"goroutines" json:"goroutines,omitempty" gorm:"column:goroutines" bson:"goroutines,omitempty" dynamodbav:"goroutines,omitempty" firestore:"goroutines,omitempty"`
}
type AuditLogSchema struct {
	Id        string    `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	User      string    `yaml:"user" mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	Ip        string    `yaml:"ip" mapstructure:"ip" json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	Resource  string    `yaml:"resource" mapstructure:"resource" json:"resource,omitempty" gorm:"column:resource" bson:"resource,omitempty" dynamodbav:"resource,omitempty" firestore:"resource,omitempty"`
	Action    string    `yaml:"action" mapstructure:"action" json:"action,omitempty" gorm:"column:action" bson:"action,omitempty" dynamodbav:"action,omitempty" firestore:"action,omitempty"`
	Timestamp string    `yaml:"timestamp" mapstructure:"timestamp" json:"timestamp,omitempty" gorm:"column:timestamp" bson:"timestamp,omitempty" dynamodbav:"timestamp,omitempty" firestore:"timestamp,omitempty"`
	Status    string    `yaml:"status" mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	Desc      string    `yaml:"desc" mapstructure:"desc" json:"desc,omitempty" gorm:"column:desc" bson:"desc,omitempty" dynamodbav:"desc,omitempty" firestore:"desc,omitempty"`
	Ext       *[]string `yaml:"ext" mapstructure:"ext" json:"ext,omitempty" gorm:"column:ext" bson:"ext,omitempty" dynamodbav:"ext,omitempty" firestore:"ext,omitempty"`
	Headers   *[]string `yaml:"headers" mapstructure:"headers" json:"headers,omitempty" gorm:"column:headers" bson:"headers,omitempty" dynamodbav:"headers,omitempty" firestore:"headers,omitempty"`
}
type AuditLogClient struct {
	Client    *http.Client
	Url       string
	Config    AuditLogConfig
	Schema    AuditLogSchema
	Generate  func(ctx context.Context) (string, error)
	Transform func(map[string]interface{}) map[string]interface{}
	Error     func(context.Context, string)
	Retries   []time.Duration
}

func (s *AuditLogClient) Write(ctx context.Context, resource string, action string, success bool, desc string) error {
	ch := s.Schema
	log := BuildLog(ctx, s.Schema, s.Config, s.Generate, s.Transform, resource, action, success, desc, ch.Ext)
	data, er0 := marshal(log)
	if er0 != nil {
		return er0
	}
	headers := BuildHeader(ctx, ch.Headers)
	if !s.Config.Goroutines {
		er3 := PostLog(ctx, s.Client, s.Url, data, headers, s.Error, s.Retries...)
		return er3
	} else {
		go PostLog(ctx, s.Client, s.Url, data, headers, s.Error, s.Retries...)
		return nil
	}
}

func BuildLog(ctx context.Context, ch AuditLogSchema, c AuditLogConfig, generate func(ctx context.Context) (string, error), transform func(map[string]interface{}) map[string]interface{}, resource string, action string, success bool, desc string, ext2 *[]string) map[string]interface{} {
	log := make(map[string]interface{})
	log[ch.Timestamp] = time.Now()
	log[ch.Resource] = resource
	log[ch.Action] = action
	if len(ch.Desc) > 0 {
		log[ch.Desc] = desc
	}
	if success {
		log[ch.Status] = c.True
	} else {
		log[ch.Status] = c.False
	}
	log[ch.User] = GetString(ctx, c.User)
	if len(ch.Ip) > 0 {
		log[ch.Ip] = GetString(ctx, c.Ip)
	}
	if len(ch.Id) > 0 && generate != nil {
		id, er0 := generate(ctx)
		if er0 == nil && len(id) > 0 {
			log[ch.Id] = id
		}
	}
	ext := BuildExt(ctx, ext2)
	if len(ext) > 0 {
		for k, v := range ext {
			log[k] = v
		}
	}
	m2 := log
	if transform != nil {
		m2 = transform(log)
		return m2
	}
	return m2
}
func PostLog(ctx context.Context, client *http.Client, url string, log []byte, headers *map[string]string, logError func(context.Context, string), retries ...time.Duration) error {
	l := len(retries)
	if l == 0 {
		_, err := DoWithClient(ctx, client, "POST", url, log, headers)
		return err
	} else {
		return PostWithRetries(ctx, client, url, log, headers, logError, retries)
	}
}
func PostWithRetries(ctx context.Context, client *http.Client, url string, log []byte, headers *map[string]string, logError func(context.Context, string), retries []time.Duration) error {
	_, er1 := DoWithClient(ctx, client, "POST", url, log, headers)
	if er1 == nil {
		return er1
	}
	i := 0
	err := retry(ctx, retries, func() (err error) {
		i = i + 1
		_, er2 := DoWithClient(ctx, client, "POST", url, log, headers)
		s := string(log)
		if logError != nil {
			if er2 != nil {
				s2 := fmt.Sprintf("Fail to send log after %d retries %s", i, s)
				logError(ctx, s2)
			} else {
				s2 := fmt.Sprintf("Send log successfully after %d retries %s", i, s)
				logError(ctx, s2)
			}
		}
		return er2
	})
	if err != nil {
		if logError != nil {
			s := string(log)
			s2 := fmt.Sprintf("Failed to send log: %s. Error: %v.", s, err)
			logError(ctx, s2)
		}
	}
	return err
}
func BuildExt(ctx context.Context, keys *[]string) map[string]interface{} {
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
func BuildHeader(ctx context.Context, keys *[]string) *map[string]string {
	if keys != nil {
		headers := make(map[string]string)
		hs := *keys
		for _, header := range hs {
			v := ctx.Value(header)
			if v != nil {
				s, ok := v.(string)
				if ok {
					headers[header] = s
				}
			}
		}
		if len(headers) > 0 {
			return &headers
		} else {
			return nil
		}
	}
	return nil
}
func GetString(ctx context.Context, key string) string {
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
func DoWithClient(ctx context.Context, client *http.Client, method string, url string, obj interface{}, headers *map[string]string) (*json.Decoder, error) {
	rq, err := marshal(obj)
	if err != nil {
		return nil, err
	}
	return DoAndBuildDecoder(ctx, client, url, method, rq, headers)
}
func DoAndBuildDecoder(ctx context.Context, client *http.Client, url string, method string, body []byte, headers *map[string]string) (*json.Decoder, error) {
	res, er1 := Do(ctx, client, url, method, body, headers)
	if er1 != nil {
		return nil, er1
	}
	if res.StatusCode == 503 {
		er2 := errors.New("503 Service Unavailable")
		return nil, er2
	}
	return json.NewDecoder(res.Body), nil
}
func Do(ctx context.Context, client *http.Client, url string, method string, body []byte, headers *map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return AddHeaderAndDo(client, req, headers)
}
func AddHeaderAndDo(client *http.Client, req *http.Request, headers *map[string]string) (*http.Response, error) {
	if headers != nil {
		for k, v := range *headers {
			req.Header.Add(k, v)
		}
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	return resp, err
}
func marshal(v interface{}) ([]byte, error) {
	b, ok1 := v.([]byte)
	if ok1 {
		return b, nil
	}
	s, ok2 := v.(string)
	if ok2 {
		return []byte(s), nil
	}
	return json.Marshal(v)
}

//Copy this code from https://stackoverflow.com/questions/47606761/repeat-code-if-an-error-occured
func retry(ctx context.Context, sleeps []time.Duration, f func() error) (err error) {
	attempts := len(sleeps)
	for i := 0; ; i++ {
		err = f()
		if err == nil {
			return
		}
		if i >= (attempts - 1) {
			break
		}
		time.Sleep(sleeps[i])
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
