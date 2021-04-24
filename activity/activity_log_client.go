package activity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func NewActivityLogClient(client *http.Client, url string, config ActivityLogConfig, schema ActivityLogSchema, generate func(context.Context) (string, error), logError func(context.Context, string), transform func(map[string]interface{}) map[string]interface{}, retries ...time.Duration) *ActivityLogClient {
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
	sender := ActivityLogClient{Client: client, Url: url, Config: config, Schema: schema, Generate: generate, Transform: transform, Error: logError, Retries: retries}
	return &sender
}

type ActivityLogConfig struct {
	User       string `mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	Ip         string `mapstructure:"ip" json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	True       string `mapstructure:"true" json:"true,omitempty" gorm:"column:true" bson:"true,omitempty" dynamodbav:"true,omitempty" firestore:"true,omitempty"`
	False      string `mapstructure:"false" json:"false,omitempty" gorm:"column:false" bson:"false,omitempty" dynamodbav:"false,omitempty" firestore:"false,omitempty"`
	Goroutines bool   `mapstructure:"goroutines" json:"goroutines,omitempty" gorm:"column:goroutines" bson:"goroutines,omitempty" dynamodbav:"goroutines,omitempty" firestore:"goroutines,omitempty"`
}
type ActivityLogSchema struct {
	Id        string    `mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	User      string    `mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	Ip        string    `mapstructure:"ip" json:"ip,omitempty" gorm:"column:ip" bson:"ip,omitempty" dynamodbav:"ip,omitempty" firestore:"ip,omitempty"`
	Resource  string    `mapstructure:"resource" json:"resource,omitempty" gorm:"column:resource" bson:"resource,omitempty" dynamodbav:"resource,omitempty" firestore:"resource,omitempty"`
	Action    string    `mapstructure:"action" json:"action,omitempty" gorm:"column:action" bson:"action,omitempty" dynamodbav:"action,omitempty" firestore:"action,omitempty"`
	Timestamp string    `mapstructure:"timestamp" json:"timestamp,omitempty" gorm:"column:timestamp" bson:"timestamp,omitempty" dynamodbav:"timestamp,omitempty" firestore:"timestamp,omitempty"`
	Status    string    `mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	Desc      string    `mapstructure:"desc" json:"desc,omitempty" gorm:"column:desc" bson:"desc,omitempty" dynamodbav:"desc,omitempty" firestore:"desc,omitempty"`
	Ext       *[]string `mapstructure:"ext" json:"ext,omitempty" gorm:"column:ext" bson:"ext,omitempty" dynamodbav:"ext,omitempty" firestore:"ext,omitempty"`
	Headers   *[]string `mapstructure:"headers" json:"headers,omitempty" gorm:"column:headers" bson:"headers,omitempty" dynamodbav:"headers,omitempty" firestore:"headers,omitempty"`
}
type ActivityLogClient struct {
	Client    *http.Client
	Url       string
	Config    ActivityLogConfig
	Schema    ActivityLogSchema
	Generate  func(ctx context.Context) (string, error)
	Transform func(map[string]interface{}) map[string]interface{}
	Error     func(context.Context, string)
	Retries   []time.Duration
}

func BuildLog(ctx context.Context, ch ActivityLogSchema, c ActivityLogConfig, generate func(ctx context.Context) (string, error), transform func(map[string]interface{}) map[string]interface{}, resource string, action string, success bool, desc string, ext2 *[]string) map[string]interface{} {
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
func (s *ActivityLogClient) Write(ctx context.Context, resource string, action string, success bool, desc string) error {
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
		//Infof(ctx, "Retrying %d of %d after error: %s", i+1, attempts, err.Error())
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
