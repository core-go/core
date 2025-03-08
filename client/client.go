package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

type ClientConfig struct {
	Endpoint Config     `yaml:"endpoint" mapstructure:"endpoint" json:"endpoint,omitempty" gorm:"column:endpoint" bson:"endpoint,omitempty" dynamodbav:"endpoint,omitempty" firestore:"endpoint,omitempty"`
	Log      *LogConfig `yaml:"log" mapstructure:"log" json:"log,omitempty" gorm:"column:log" bson:"log,omitempty" dynamodbav:"log,omitempty" firestore:"log,omitempty"`
}
type ClientConf struct {
	Config   Conf       `yaml:"config" mapstructure:"config" json:"config,omitempty" gorm:"column:config" bson:"config,omitempty" dynamodbav:"config,omitempty" firestore:"config,omitempty"`
	Endpoint Endpoint   `yaml:"endpoint" mapstructure:"endpoint" json:"endpoint,omitempty" gorm:"column:endpoint" bson:"endpoint,omitempty" dynamodbav:"endpoint,omitempty" firestore:"endpoint,omitempty"`
	Log      *LogConfig `yaml:"log" mapstructure:"log" json:"log,omitempty" gorm:"column:log" bson:"log,omitempty" dynamodbav:"log,omitempty" firestore:"log,omitempty"`
}
type Endpoint struct {
	Url      string  `yaml:"url" mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
	Username *string `yaml:"username" mapstructure:"username" json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Password *string `yaml:"password" mapstructure:"password" json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
}
type Config struct {
	Insecure *bool          `yaml:"insecure" mapstructure:"insecure" json:"insecure,omitempty" gorm:"column:insecure" bson:"insecure,omitempty" dynamodbav:"insecure,omitempty" firestore:"insecure,omitempty"`
	Timeout  *time.Duration `yaml:"timeout" mapstructure:"timeout" json:"timeout,omitempty" gorm:"column:timeout" bson:"timeout,omitempty" dynamodbav:"timeout,omitempty" firestore:"timeout,omitempty"`
	CertFile string         `yaml:"cert_file" mapstructure:"cert_file" json:"certFile,omitempty" gorm:"column:certfile" bson:"certFile,omitempty" dynamodbav:"certFile,omitempty" firestore:"certFile,omitempty"`
	KeyFile  string         `yaml:"key_file" mapstructure:"key_file" json:"keyFile,omitempty" gorm:"column:keyfile" bson:"keyFile,omitempty" dynamodbav:"keyFile,omitempty" firestore:"keyFile,omitempty"`
	PEMFile  bool           `yaml:"pem_file" mapstructure:"pem_file" json:"pemFile,omitempty" gorm:"column:pemFile" bson:"pemFile,omitempty" dynamodbav:"pemFile,omitempty" firestore:"pemFile,omitempty"`
	Url      string         `yaml:"url" mapstructure:"url" json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty"`
	Username *string        `yaml:"username" mapstructure:"username" json:"username,omitempty" gorm:"column:username" bson:"username,omitempty" dynamodbav:"username,omitempty" firestore:"username,omitempty"`
	Password *string        `yaml:"password" mapstructure:"password" json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
}
type Conf struct {
	Insecure *bool          `yaml:"insecure" mapstructure:"insecure" json:"insecure,omitempty" gorm:"column:insecure" bson:"insecure,omitempty" dynamodbav:"insecure,omitempty" firestore:"insecure,omitempty"`
	Timeout  *time.Duration `yaml:"timeout" mapstructure:"timeout" json:"timeout,omitempty" gorm:"column:timeout" bson:"timeout,omitempty" dynamodbav:"timeout,omitempty" firestore:"timeout,omitempty"`
	CertFile string         `yaml:"cert_file" mapstructure:"cert_file" json:"certFile,omitempty" gorm:"column:certfile" bson:"certFile,omitempty" dynamodbav:"certFile,omitempty" firestore:"certFile,omitempty"`
	KeyFile  string         `yaml:"key_file" mapstructure:"key_file" json:"keyFile,omitempty" gorm:"column:keyfile" bson:"keyFile,omitempty" dynamodbav:"keyFile,omitempty" firestore:"keyFile,omitempty"`
	PEMFile  bool           `yaml:"pem_file" mapstructure:"pem_file" json:"pemFile,omitempty" gorm:"column:pemFile" bson:"pemFile,omitempty" dynamodbav:"pemFile,omitempty" firestore:"pemFile,omitempty"`
}
type LogConfig struct {
	Separate       bool   `yaml:"separate" mapstructure:"separate" json:"separate,omitempty" gorm:"column:separate" bson:"separate,omitempty" dynamodbav:"separate,omitempty" firestore:"separate,omitempty"`
	Log            bool   `yaml:"log" mapstructure:"log" json:"log,omitempty" gorm:"column:log" bson:"log,omitempty" dynamodbav:"log,omitempty" firestore:"log,omitempty"`
	Duration       string `yaml:"duration" mapstructure:"duration" json:"duration,omitempty" gorm:"column:duration" bson:"duration,omitempty" dynamodbav:"duration,omitempty" firestore:"duration,omitempty"`
	Size           string `yaml:"size" mapstructure:"size" json:"size,omitempty" gorm:"column:size" bson:"size,omitempty" dynamodbav:"size,omitempty" firestore:"size,omitempty"`
	ResponseStatus string `yaml:"status" mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	Request        string `yaml:"request" mapstructure:"request" json:"request,omitempty" gorm:"column:request" bson:"request,omitempty" dynamodbav:"request,omitempty" firestore:"request,omitempty"`
	Response       string `yaml:"response" mapstructure:"response" json:"response,omitempty" gorm:"column:response" bson:"response,omitempty" dynamodbav:"response,omitempty" firestore:"response,omitempty"`
	Error          string `yaml:"error" mapstructure:"error" json:"error,omitempty" gorm:"column:error" bson:"error,omitempty" dynamodbav:"error,omitempty" firestore:"error,omitempty"`
}
type Params struct {
	Client   *http.Client
	Url      string
	Header   map[string]string
	Config   *LogConfig
	LogError func(context.Context, string, map[string]interface{})
	LogInfo  func(context.Context, string, map[string]interface{})
}

const (
	post   = "POST"
	put    = "PUT"
	get    = "GET"
	patch  = "PATCH"
	delete = "DELETE"
)

// var conf3 LogConfig
var sClient *http.Client

func SetClient(c *http.Client) {
	sClient = c
}
func InitializeLog(c *LogConfig) *LogConfig {
	var c2 LogConfig
	if c == nil {
		c2.Log = true
		c2.ResponseStatus = "status"
		c2.Size = "size"
		c2.Duration = "duration"
		c2.Error = "error"
		return &c2
	}
	c2.Log = c.Log
	c2.Separate = c.Separate
	c2.ResponseStatus = c.ResponseStatus
	c2.Size = c.Size
	if len(c.Duration) > 0 {
		c2.Duration = c.Duration
	} else {
		c2.Duration = "duration"
	}
	if len(c.Error) > 0 {
		c2.Error = c.Error
	} else {
		c2.Error = "error"
	}
	c2.Request = c.Request
	c2.Response = c.Response
	return &c2
}
func InitializeParams(config ClientConfig, opts ...func(context.Context, string, map[string]interface{})) (*Params, error) {
	c, header, conf, err := InitializeClient(config)
	if err != nil {
		return nil, err
	}
	var logError func(context.Context, string, map[string]interface{})
	var logInfo func(context.Context, string, map[string]interface{})
	if len(opts) > 0 && opts[0] != nil {
		logError = opts[0]
	}
	if len(opts) > 1 && opts[1] != nil {
		logInfo = opts[1]
	}
	return &Params{Client: c, Url: config.Endpoint.Url, Header: header, Config: conf, LogError: logError, LogInfo: logInfo}, nil
}
func InitParams(config ClientConf, opts ...func(context.Context, string, map[string]interface{})) (*Params, error) {
	c, header, conf, err := InitClient(config)
	if err != nil {
		return nil, err
	}
	var logError func(context.Context, string, map[string]interface{})
	var logInfo func(context.Context, string, map[string]interface{})
	if len(opts) > 0 && opts[0] != nil {
		logError = opts[0]
	}
	if len(opts) > 1 && opts[1] != nil {
		logInfo = opts[1]
	}
	return &Params{Client: c, Url: config.Endpoint.Url, Header: header, Config: conf, LogError: logError, LogInfo: logInfo}, nil
}
func InitializeClient(config ClientConfig) (*http.Client, map[string]string, *LogConfig, error) {
	e := config.Endpoint
	conf := Conf{
		Insecure: e.Insecure,
		Timeout:  e.Timeout,
		CertFile: e.CertFile,
		KeyFile:  e.KeyFile,
		PEMFile:  e.PEMFile,
	}
	c, err := NewClient(conf)
	if err != nil {
		return nil, nil, nil, err
	}
	header := CreateHeaderFromConfig(config.Endpoint)
	l := InitializeLog(config.Log)
	return c, header, l, nil
}
func InitClient(config ClientConf) (*http.Client, map[string]string, *LogConfig, error) {
	c, err := NewClient(config.Config)
	if err != nil {
		return nil, nil, nil, err
	}
	header := CreateHeaderFromConf(config.Endpoint)
	l := InitializeLog(config.Log)
	return c, header, l, nil
}
func NewClient(c Conf) (*http.Client, error) {
	if len(c.CertFile) > 0 && len(c.KeyFile) > 0 {
		return NewTLSClient(c.CertFile, c.KeyFile, c.Timeout)
	} else {
		if c.Insecure != nil {
			if c.Timeout != nil {
				transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: *c.Insecure}}
				client0 := &http.Client{Transport: transport, Timeout: *c.Timeout}
				// sClient = client0
				return client0, nil
			} else {
				transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: *c.Insecure}}
				client0 := &http.Client{Transport: transport}
				return client0, nil
			}
		} else {
			if c.Timeout != nil {
				client0 := &http.Client{Timeout: *c.Timeout}
				return client0, nil
			} else {
				client0 := &http.Client{}
				return client0, nil
			}
		}
	}
}
func NewTLSClient(certFile, keyFile string, timeout *time.Duration, options ...string) (*http.Client, error) {
	clientCert, er1 := tls.LoadX509KeyPair(certFile, keyFile)
	if er1 != nil {
		return nil, er1
	}
	conf, er2 := GetTLSClientConfig(clientCert, options...)
	if er2 != nil {
		return nil, er2
	}
	if timeout != nil {
		client0 := &http.Client{Transport: &http.Transport{TLSClientConfig: conf}}
		// sClient = client0
		return client0, nil
	} else {
		client0 := &http.Client{
			Transport: &http.Transport{TLSClientConfig: conf},
			Timeout:   *timeout,
		}
		// sClient = client0
		return client0, nil
	}
}
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
func CreateHeader(username, password string) map[string]string {
	h := make(map[string]string, 0)
	h["Authorization"] = "Basic " + BasicAuth(username, password)
	return h
}
func CreateHeaderFromConf(c Endpoint) map[string]string {
	if c.Username == nil || c.Password == nil || len(*c.Username) == 0 {
		return nil
	}
	h := make(map[string]string, 0)
	h["Authorization"] = "Basic " + BasicAuth(*c.Username, *c.Password)
	return h
}
func CreateHeaderFromConfig(c Config) map[string]string {
	if c.Username == nil || c.Password == nil || len(*c.Username) == 0 {
		return nil
	}
	h := make(map[string]string, 0)
	h["Authorization"] = "Basic " + BasicAuth(*c.Username, *c.Password)
	return h
}
func GetTLSClientConfig(clientCert tls.Certificate, options ...string) (*tls.Config, error) {
	c := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{clientCert},
		MinVersion:         tls.VersionTLS10,
		MaxVersion:         tls.VersionTLS13,
	}
	if len(options) > 0 && len(options[0]) > 0 {
		pem, err := os.ReadFile(options[0])
		if err != nil {
			return nil, err
		}
		roots := x509.NewCertPool()
		roots.AppendCertsFromPEM(pem)
		c.RootCAs = roots
	}
	return c, nil
}
func DoRequest(ctx context.Context, client *http.Client, method string, url string, body []byte, headers map[string]string) (*http.Response, error) {
	if body != nil {
		b := body
		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(b))
		if err != nil {
			return nil, err
		}
		return AddHeaderAndDo(client, req, headers)
	} else {
		req, err := http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, err
		}
		return AddHeaderAndDo(client, req, headers)
	}
}
func DoJSON(ctx context.Context, client *http.Client, method string, url string, body []byte, headers map[string]string) (*http.Response, error) {
	if body != nil {
		b := body
		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(b))
		if err != nil {
			return nil, err
		}
		return AddHeaderAndDoJSON(client, req, headers)
	} else {
		req, err := http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, err
		}
		return AddHeaderAndDoJSON(client, req, headers)
	}
}
func DoJSONWithClient(ctx context.Context, client *http.Client, method string, url string, obj interface{}, headers map[string]string, errorStatus int) (*json.Decoder, error) {
	if client == nil {
		client = sClient
	}
	rq, err := Marshal(obj)
	if err != nil {
		return nil, err
	}
	return DoJSONAndDecode(ctx, client, method, url, rq, headers, errorStatus)
}
func DoJSONAndDecode(ctx context.Context, client *http.Client, method string, url string, body []byte, headers map[string]string, errorStatus int) (*json.Decoder, error) {
	start := time.Now()
	response, er1 := DoJSON(ctx, client, method, url, body, headers)
	end := time.Now()
	dur := end.Sub(start).Milliseconds()
	if er1 != nil {
		return nil, er1
	}
	if errorStatus < 0 {
		res := json.NewDecoder(response.Body)
		return res, nil
	}
	if response.StatusCode >= errorStatus {
		res, er2 := io.ReadAll(response.Body)
		var rs string
		if er2 == nil {
			rs = string(res)
		}
		return nil, NewHttpError(response.StatusCode, nil, dur, fmt.Sprint("Response error with status code: ", response.StatusCode), url, string(body), rs)
	}
	res := json.NewDecoder(response.Body)
	return res, nil
}
func AddHeaderAndDoJSON(client *http.Client, req *http.Request, headers map[string]string) (*http.Response, error) {
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	return resp, err
}
func AddHeaderAndDo(client *http.Client, req *http.Request, headers map[string]string) (*http.Response, error) {
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	return resp, err
}
func DoGet(ctx context.Context, client *http.Client, url string, headers map[string]string) (*http.Response, error) {
	return DoJSON(ctx, client, get, url, nil, headers)
}
func DoDelete(ctx context.Context, client *http.Client, url string, headers map[string]string) (*http.Response, error) {
	return DoJSON(ctx, client, delete, url, nil, headers)
}
func DoPost(ctx context.Context, client *http.Client, url string, body []byte, headers map[string]string) (*http.Response, error) {
	return DoJSON(ctx, client, post, url, body, headers)
}
func DoPut(ctx context.Context, client *http.Client, url string, body []byte, headers map[string]string) (*http.Response, error) {
	return DoJSON(ctx, client, put, url, body, headers)
}
func DoPatch(ctx context.Context, client *http.Client, url string, body []byte, headers map[string]string) (*http.Response, error) {
	return DoJSON(ctx, client, patch, url, body, headers)
}
func GetDecoder(ctx context.Context, client *http.Client, url string, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, client, get, url, nil, nil, conf, options...)
}
func GetDecoderWithHeader(ctx context.Context, client *http.Client, url string, headers map[string]string, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, client, get, url, nil, headers, conf, options...)
}
func Get(ctx context.Context, client *http.Client, url string, result interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) error {
	return GetWithHeader(ctx, client, url, nil, result, conf, options...)
}
func GetWithHeader(ctx context.Context, client *http.Client, url string, headers map[string]string, result interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) error {
	decoder, er1 := DoWithClient(ctx, client, get, url, nil, headers, conf, options...)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func DeleteDecoder(ctx context.Context, client *http.Client, url string, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, client, delete, url, nil, nil, conf, options...)
}
func DeleteDecoderWithHeader(ctx context.Context, client *http.Client, url string, headers map[string]string, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, client, delete, url, nil, headers, conf, options...)
}
func Delete(ctx context.Context, client *http.Client, url string, result interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) error {
	return DeleteWithHeader(ctx, client, url, nil, result, conf, options...)
}
func DeleteWithHeader(ctx context.Context, client *http.Client, url string, headers map[string]string, result interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) error {
	decoder, er1 := DoWithClient(ctx, client, delete, url, nil, headers, conf, options...)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func PostDecoder(ctx context.Context, client *http.Client, url string, obj interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, client, post, url, obj, nil, conf, options...)
}
func PostDecoderWithHeader(ctx context.Context, client *http.Client, url string, obj interface{}, headers map[string]string, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, client, post, url, obj, headers, conf, options...)
}
func Post(ctx context.Context, client *http.Client, url string, obj interface{}, result interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) error {
	return PostWithHeader(ctx, client, url, obj, nil, result, conf, options...)
}
func PostWithHeader(ctx context.Context, client *http.Client, url string, obj interface{}, headers map[string]string, result interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) error {
	decoder, er1 := DoWithClient(ctx, client, post, url, obj, headers, conf, options...)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func PutDecoder(ctx context.Context, client *http.Client, url string, obj interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, client, put, url, obj, nil, conf, options...)
}
func PutDecoderWithHeader(ctx context.Context, client *http.Client, url string, obj interface{}, headers map[string]string, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, client, put, url, obj, headers, conf, options...)
}
func Put(ctx context.Context, client *http.Client, url string, obj interface{}, result interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) error {
	return PutWithHeader(ctx, client, url, obj, nil, result, conf, options...)
}
func PutWithHeader(ctx context.Context, client *http.Client, url string, obj interface{}, headers map[string]string, result interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) error {
	decoder, er1 := DoWithClient(ctx, client, put, url, obj, headers, conf, options...)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func PatchDecoder(ctx context.Context, client *http.Client, url string, obj interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, client, patch, url, obj, nil, conf, options...)
}
func PatchDecoderWithHeader(ctx context.Context, client *http.Client, url string, obj interface{}, headers map[string]string, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, client, patch, url, obj, headers, conf, options...)
}
func Patch(ctx context.Context, client *http.Client, url string, obj interface{}, result interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) error {
	return PatchWithHeader(ctx, client, url, obj, nil, result, conf, options...)
}
func PatchWithHeader(ctx context.Context, client *http.Client, url string, obj interface{}, headers map[string]string, result interface{}, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) error {
	decoder, er1 := DoWithClient(ctx, client, patch, url, obj, headers, conf, options...)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func Marshal(obj interface{}) ([]byte, error) {
	b, ok := obj.([]byte)
	if ok {
		return b, nil
	}
	s, ok2 := obj.(string)
	if ok2 {
		b2 := []byte(s)
		return b2, nil
	}
	v, er0 := json.Marshal(obj)
	if er0 != nil {
		return nil, er0
	}
	return v, nil
}
func GetString(obj interface{}) (string, bool) {
	bs, err := Marshal(obj)
	if err != nil {
		return "", false
	}
	return string(bs), true
}
func DoWithClient(ctx context.Context, client *http.Client, method string, url string, obj interface{}, headers map[string]string, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	if client == nil {
		client = sClient
	}
	rq, err := Marshal(obj)
	if err != nil {
		return nil, err
	}
	return DoAndBuildDecoder(ctx, client, method, url, rq, headers, conf, options...)
}
func DoAndBuildDecoder(ctx context.Context, client *http.Client, method string, url string, body []byte, headers map[string]string, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	var logError func(context.Context, string, map[string]interface{})
	var logInfo func(context.Context, string, map[string]interface{})
	if len(options) > 0 {
		logError = options[0]
	}
	if len(options) > 1 {
		logInfo = options[1]
	}
	start := time.Now()
	res, er1 := DoJSON(ctx, client, method, url, body, headers)
	end := time.Now()
	dur := end.Sub(start).Milliseconds()
	if logError != nil && (er1 != nil || res.StatusCode >= 400) {
		fs3 := make(map[string]interface{}, 0)
		var c2 LogConfig
		if conf != nil {
			c2 = *conf
		} else {
			c2.Duration = "duration"
			c2.Request = "request"
			c2.Response = "response"
			c2.ResponseStatus = "status"
			c2.Error = "error"
		}
		if body != nil {
			rq := string(body)
			if len(rq) > 0 {
				fs3[c2.Request] = rq
			}
		}
		if er1 != nil {
			if len(c2.Error) > 0 {
				fs3[c2.Error] = er1.Error()
			}
			logError(ctx, method+" "+url, fs3)
			return nil, er1
		}
		if len(c2.ResponseStatus) > 0 {
			fs3[c2.ResponseStatus] = res.StatusCode
		}
		if len(c2.Size) > 0 && res.ContentLength > 0 {
			fs3[c2.Size] = res.ContentLength
		}
		if len(c2.Response) > 0 {
			dump, er3 := httputil.DumpResponse(res, true)
			if er3 != nil {
				if len(c2.Error) > 0 {
					fs3[c2.Error] = er3.Error()
				}
				logError(ctx, method+" "+url, fs3)
				return nil, er3
			}
			s := string(dump)
			if len(c2.Size) > 0 {
				fs3[c2.Size] = len(s)
			}
			if len(c2.Response) > 0 {
				fs3[c2.Response] = s
			}
			if res.StatusCode == 503 {
				logError(ctx, method+" "+url, fs3)
				var rq string
				if body != nil {
					rq = string(body)
				}
				er2 := NewHttpError(http.StatusServiceUnavailable, er1, dur, "503 Service Unavailable", url, rq, s)
				return nil, er2
			}
			logError(ctx, method+" "+url, fs3)
			return json.NewDecoder(strings.NewReader(s)), nil
		} else {
			if res.StatusCode == 503 {
				logError(ctx, method+" "+url, fs3)
				var rq string
				if body != nil {
					rq = string(body)
				}
				er2 := NewHttpError(http.StatusServiceUnavailable, er1, dur, "503 Service Unavailable", url, rq)
				return nil, er2
			}
			logError(ctx, method+" "+url, fs3)
			return json.NewDecoder(res.Body), nil
		}
	}
	if conf != nil && conf.Log == true && logInfo != nil {
		canRequest := false
		if method != "GET" && method != "DELETE" && method != "OPTIONS" {
			canRequest = true
		}
		if conf.Separate && len(conf.Request) > 0 && body != nil && canRequest {
			fs1 := make(map[string]interface{}, 0)
			rq := string(body)
			if len(rq) > 0 {
				fs1[conf.Request] = rq
			}
			logInfo(ctx, method+" "+url, fs1)
		}
		fs3 := make(map[string]interface{}, 0)
		fs3[conf.Duration] = dur
		if !conf.Separate && len(conf.Request) > 0 && body != nil && canRequest {
			rq := string(body)
			if len(rq) > 0 {
				fs3[conf.Request] = rq
			}
		}
		if er1 != nil {
			if len(conf.Error) > 0 {
				fs3[conf.Error] = er1.Error()
			}
			logInfo(ctx, method+" "+url, fs3)
			return nil, er1
		}
		if len(conf.ResponseStatus) > 0 {
			fs3[conf.ResponseStatus] = res.StatusCode
		}
		if len(conf.Size) > 0 && res.ContentLength > 0 {
			fs3[conf.Size] = res.ContentLength
		}
		if len(conf.Response) > 0 {
			buf := new(bytes.Buffer)
			_, er3 := buf.ReadFrom(res.Body)
			if er3 != nil {
				if len(conf.Error) > 0 {
					fs3[conf.Error] = er3.Error()
				}
				logInfo(ctx, method+" "+url, fs3)
				return nil, er3
			}
			s := buf.String()
			if len(conf.Size) > 0 {
				fs3[conf.Size] = len(s)
			}
			if len(conf.Response) > 0 {
				fs3[conf.Response] = s
			}
			if res.StatusCode == 503 {
				logInfo(ctx, method+" "+url, fs3)
				var rq string
				if body != nil {
					rq = string(body)
				}
				er2 := NewHttpError(http.StatusServiceUnavailable, er1, dur, "503 Service Unavailable", url, rq)
				return nil, er2
			}
			logInfo(ctx, method+" "+url, fs3)
			return json.NewDecoder(strings.NewReader(s)), nil
		} else {
			if res.StatusCode == 503 {
				logInfo(ctx, method+" "+url, fs3)
				var rq string
				if body != nil {
					rq = string(body)
				}
				er2 := NewHttpError(http.StatusServiceUnavailable, er1, dur, "503 Service Unavailable", url, rq)
				return nil, er2
			}
			logInfo(ctx, method+" "+url, fs3)
			return json.NewDecoder(res.Body), nil
		}
	} else {
		if er1 != nil {
			return nil, er1
		}
		if res.StatusCode == 503 {
			var rq string
			if body != nil {
				rq = string(body)
			}
			er2 := NewHttpError(http.StatusServiceUnavailable, er1, dur, "503 Service Unavailable", url, rq)
			return nil, er2
		}
		return json.NewDecoder(res.Body), nil
	}
}

func DoAndLog(ctx context.Context, client *http.Client, method string, url string, body []byte, headers map[string]string, conf *LogConfig, options ...func(context.Context, string, map[string]interface{})) (*http.Response, error) {
	var logError func(context.Context, string, map[string]interface{})
	var logInfo func(context.Context, string, map[string]interface{})
	if len(options) > 0 {
		logError = options[0]
	}
	if len(options) > 1 {
		logInfo = options[1]
	}
	start := time.Now()
	res, er1 := DoJSON(ctx, client, method, url, body, headers)
	end := time.Now()
	dur := end.Sub(start).Milliseconds()
	if logError != nil && (er1 != nil || res.StatusCode >= 400) {
		fs3 := make(map[string]interface{}, 0)
		var c2 LogConfig
		if conf != nil {
			c2 = *conf
		} else {
			c2.Duration = "duration"
			c2.Request = "request"
			c2.Response = "response"
			c2.ResponseStatus = "status"
			c2.Error = "error"
		}
		if body != nil {
			rq := string(body)
			if len(rq) > 0 {
				fs3[c2.Request] = rq
			}
		}
		if er1 != nil {
			if len(c2.Error) > 0 {
				fs3[c2.Error] = er1.Error()
			}
			logError(ctx, method+" "+url, fs3)
			return res, er1
		}
		if len(c2.ResponseStatus) > 0 {
			fs3[c2.ResponseStatus] = res.StatusCode
		}
		if len(c2.Size) > 0 && res.ContentLength > 0 {
			fs3[c2.Size] = res.ContentLength
		}
		if len(c2.Response) > 0 {
			dump, er3 := httputil.DumpResponse(res, true)
			if er3 != nil {
				if len(c2.Error) > 0 {
					fs3[c2.Error] = er3.Error()
				}
				logError(ctx, method+" "+url, fs3)
				return res, er3
			}
			s := string(dump)
			if len(c2.Size) > 0 {
				fs3[c2.Size] = len(s)
			}
			if len(c2.Response) > 0 {
				fs3[c2.Response] = s
			}
			if res.StatusCode == 503 {
				logError(ctx, method+" "+url, fs3)
				var rq string
				if body != nil {
					rq = string(body)
				}
				er2 := NewHttpError(http.StatusServiceUnavailable, er1, dur, "503 Service Unavailable", url, rq, s)
				return res, er2
			}
			logError(ctx, method+" "+url, fs3)
			return res, nil
		} else {
			if res.StatusCode == 503 {
				logError(ctx, method+" "+url, fs3)
				var rq string
				if body != nil {
					rq = string(body)
				}
				er2 := NewHttpError(http.StatusServiceUnavailable, er1, dur, "503 Service Unavailable", url, rq)
				return res, er2
			}
			logError(ctx, method+" "+url, fs3)
			return res, nil
		}
	}
	if conf != nil && conf.Log == true && logInfo != nil {
		canRequest := false
		if method != "GET" && method != "DELETE" && method != "OPTIONS" {
			canRequest = true
		}
		if conf.Separate && len(conf.Request) > 0 && body != nil && canRequest {
			fs1 := make(map[string]interface{}, 0)
			rq := string(body)
			if len(rq) > 0 {
				fs1[conf.Request] = rq
			}
			logInfo(ctx, method+" "+url, fs1)
		}
		fs3 := make(map[string]interface{}, 0)
		fs3[conf.Duration] = dur
		if !conf.Separate && len(conf.Request) > 0 && body != nil && canRequest {
			rq := string(body)
			if len(rq) > 0 {
				fs3[conf.Request] = rq
			}
		}
		if er1 != nil {
			if len(conf.Error) > 0 {
				fs3[conf.Error] = er1.Error()
			}
			logInfo(ctx, method+" "+url, fs3)
			return res, er1
		}
		if len(conf.ResponseStatus) > 0 {
			fs3[conf.ResponseStatus] = res.StatusCode
		}
		if len(conf.Size) > 0 && res.ContentLength > 0 {
			fs3[conf.Size] = res.ContentLength
		}
		if len(conf.Response) > 0 {
			dump, er3 := httputil.DumpResponse(res, true)
			if er3 != nil {
				if len(conf.Error) > 0 {
					fs3[conf.Error] = er3.Error()
				}
				logInfo(ctx, method+" "+url, fs3)
				return res, er3
			}
			s := string(dump)
			if len(conf.Size) > 0 {
				fs3[conf.Size] = len(s)
			}
			if len(conf.Response) > 0 {
				fs3[conf.Response] = s
			}
			if res.StatusCode == 503 {
				logInfo(ctx, method+" "+url, fs3)
				var rq string
				if body != nil {
					rq = string(body)
				}
				er2 := NewHttpError(http.StatusServiceUnavailable, er1, dur, "503 Service Unavailable", url, rq, s)
				return res, er2
			}
			logInfo(ctx, method+" "+url, fs3)
			return res, nil
		} else {
			if res.StatusCode == 503 {
				logInfo(ctx, method+" "+url, fs3)
				var rq string
				if body != nil {
					rq = string(body)
				}
				er2 := NewHttpError(http.StatusServiceUnavailable, er1, dur, "503 Service Unavailable", url, rq)
				return res, er2
			}
			logInfo(ctx, method+" "+url, fs3)
			return res, nil
		}
	} else {
		if er1 != nil {
			return nil, er1
		}
		if res.StatusCode == 503 {
			var rq string
			if body != nil {
				rq = string(body)
			}
			er2 := NewHttpError(http.StatusServiceUnavailable, er1, dur, "503 Service Unavailable", url, rq)
			return res, er2
		}
		return res, nil
	}
}

type HttpError struct {
	StatusCode   int
	ErrorMessage string
	RootError    error
	Url          string
	Request      string
	Response     string
	Duration     int64
	ErrorType    string
	ErrorCode    string
	Service      string
	Severity     string
}

func NewHttpError(statusCode int, rootError error, duration int64, opts ...string) error {
	err := &HttpError{StatusCode: statusCode, RootError: rootError, Duration: duration}
	if len(opts) > 0 {
		err.ErrorMessage = opts[0]
	} else if rootError != nil {
		err.ErrorMessage = rootError.Error()
	}
	if len(opts) > 1 {
		err.Url = opts[1]
	}
	if len(opts) > 2 {
		err.Request = opts[2]
	}
	if len(opts) > 3 {
		err.Response = opts[3]
	}
	if len(opts) > 4 {
		err.ErrorType = opts[4]
	}
	if len(opts) > 5 {
		err.ErrorCode = opts[5]
	}
	if len(opts) > 6 {
		err.Service = opts[6]
	}
	if len(opts) > 7 {
		err.Severity = opts[7]
	}
	return err
}
func (e *HttpError) Error() string {
	if len(e.ErrorMessage) > 0 {
		return e.ErrorMessage
	}
	return e.GetRootError()
}
func (e *HttpError) GetRootError() string {
	if e.RootError != nil {
		return e.RootError.Error()
	}
	return ""
}
func IsHttpError(err error) (*HttpError, bool) {
	httpErr, ok := err.(*HttpError)
	if ok {
		return httpErr, ok
	} else {
		return nil, ok
	}
}
func MakeMap(err *HttpError, prefix string) map[string]interface{} {
	mp := make(map[string]interface{})
	mp[prefix+"Duration"] = err.Duration
	mp[prefix+"Status"] = err.StatusCode
	if len(err.Request) > 0 {
		mp[prefix+"Request"] = err.Request
	}
	if len(err.Response) > 0 {
		mp[prefix+"Response"] = err.Response
	}
	if len(err.Url) > 0 {
		mp[prefix+"Url"] = err.Url
	}
	if len(err.ErrorMessage) > 0 {
		mp[prefix+"Error"] = err.ErrorMessage
	}
	if err.RootError != nil && err.Error() != err.ErrorMessage {
		mp[prefix+"RootError"] = err.Error()
	}
	if len(err.ErrorType) > 0 {
		mp[prefix+"ErrorType"] = err.ErrorType
	}
	if len(err.ErrorCode) > 0 {
		mp[prefix+"ErrorCode"] = err.ErrorCode
	}
	if len(err.Service) > 0 {
		mp[prefix+"Service"] = err.Service
	}
	if len(err.Severity) > 0 {
		mp[prefix+"Severity"] = err.Severity
	}
	return mp
}
