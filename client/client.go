package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	Insecure       *bool  `mapstructure:"insecure" json:"insecure,omitempty" gorm:"column:insecure" bson:"insecure,omitempty" dynamodbav:"insecure,omitempty" firestore:"insecure,omitempty"`
	Timeout        int64  `mapstructure:"timeout" json:"timeout,omitempty" gorm:"column:timeout" bson:"timeout,omitempty" dynamodbav:"timeout,omitempty" firestore:"timeout,omitempty"`
	CertFile       string `mapstructure:"cert_file" json:"certFile,omitempty" gorm:"column:certfile" bson:"certFile,omitempty" dynamodbav:"certFile,omitempty" firestore:"certFile,omitempty"`
	KeyFile        string `mapstructure:"key_file" json:"keyFile,omitempty" gorm:"column:keyfile" bson:"keyFile,omitempty" dynamodbav:"keyFile,omitempty" firestore:"keyFile,omitempty"`
	Separate       bool   `mapstructure:"separate" json:"separate,omitempty" gorm:"column:separate" bson:"separate,omitempty" dynamodbav:"separate,omitempty" firestore:"separate,omitempty"`
	Log            bool   `mapstructure:"log" json:"log,omitempty" gorm:"column:log" bson:"log,omitempty" dynamodbav:"log,omitempty" firestore:"log,omitempty"`
	Duration       string `mapstructure:"duration" json:"duration,omitempty" gorm:"column:duration" bson:"duration,omitempty" dynamodbav:"duration,omitempty" firestore:"duration,omitempty"`
	Size           string `mapstructure:"size" json:"size,omitempty" gorm:"column:size" bson:"size,omitempty" dynamodbav:"size,omitempty" firestore:"size,omitempty"`
	ResponseStatus string `mapstructure:"status" json:"status,omitempty" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	Request        string `mapstructure:"request" json:"request,omitempty" gorm:"column:request" bson:"request,omitempty" dynamodbav:"request,omitempty" firestore:"request,omitempty"`
	Response       string `mapstructure:"response" json:"response,omitempty" gorm:"column:response" bson:"response,omitempty" dynamodbav:"response,omitempty" firestore:"response,omitempty"`
	Error          string `mapstructure:"error" json:"error,omitempty" gorm:"column:error" bson:"error,omitempty" dynamodbav:"error,omitempty" firestore:"error,omitempty"`
}

const (
	post   = "POST"
	put    = "PUT"
	get    = "GET"
	patch  = "PATCH"
	delete = "DELETE"
)

var conf Config
var sClient *http.Client

func SetClient(c *http.Client) {
	sClient = c
}

func NewClient(c Config) (*http.Client, error) {
	conf.Log = c.Log
	conf.Separate = c.Separate
	conf.ResponseStatus = c.ResponseStatus
	conf.Size = c.Size
	if len(c.Duration) > 0 {
		conf.Duration = c.Duration
	} else {
		conf.Duration = "duration"
	}
	if len(c.Request) > 0 {
		conf.Request = c.Request
	} else {
		conf.Request = "request"
	}
	if len(c.Response) > 0 {
		conf.Response = c.Response
	} else {
		conf.Response = "response"
	}
	if len(c.Error) > 0 {
		conf.Error = c.Error
	} else {
		conf.Error = "error"
	}
	if len(c.CertFile) > 0 && len(c.KeyFile) > 0 {
		return NewTLSClient(c.CertFile, c.KeyFile, time.Duration(c.Timeout)*time.Millisecond)
	} else {
		if c.Insecure != nil {
			if c.Timeout > 0 {
				transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: *c.Insecure}}
				client0 := &http.Client{Transport: transport, Timeout: time.Duration(c.Timeout) * time.Millisecond}
				sClient = client0
				return client0, nil
			} else {
				transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: *c.Insecure}}
				client0 := &http.Client{Transport: transport}
				sClient = client0
				return client0, nil
			}
		} else {
			if c.Timeout > 0 {
				client0 := &http.Client{Timeout: time.Duration(c.Timeout) * time.Millisecond}
				sClient = client0
				return client0, nil
			} else {
				client0 := &http.Client{}
				sClient = client0
				return client0, nil
			}
		}
	}
}
func NewTLSClient(certFile, keyFile string, timeout time.Duration) (*http.Client, error) {
	clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	if timeout <= 0 {
		client0 := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					Certificates:       []tls.Certificate{clientCert},
					MinVersion:         tls.VersionTLS10,
					MaxVersion:         tls.VersionTLS10,
				},
			},
		}
		sClient = client0
		return client0, nil
	} else {
		client0 := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					Certificates:       []tls.Certificate{clientCert},
					MinVersion:         tls.VersionTLS10,
					MaxVersion:         tls.VersionTLS10,
				},
			},
			Timeout: timeout * time.Second,
		}
		sClient = client0
		return client0, nil
	}
}
func DoJSON(ctx context.Context, client *http.Client, url string, method string, body []byte, headers map[string]string) (*http.Response, error) {
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
	return DoJSON(ctx, client, url, get, nil, headers)
}
func DoDelete(ctx context.Context, client *http.Client, url string, headers map[string]string) (*http.Response, error) {
	return DoJSON(ctx, client, url, delete, nil, headers)
}
func DoPost(ctx context.Context, client *http.Client, url string, body []byte, headers map[string]string) (*http.Response, error) {
	return DoJSON(ctx, client, url, post, body, headers)
}
func DoPut(ctx context.Context, client *http.Client, url string, body []byte, headers map[string]string) (*http.Response, error) {
	return DoJSON(ctx, client, url, put, body, headers)
}
func DoPatch(ctx context.Context, client *http.Client, url string, body []byte, headers map[string]string) (*http.Response, error) {
	return DoJSON(ctx, client, url, patch, body, headers)
}
func Get(ctx context.Context, url string, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, sClient, get, url, nil, nil, options...)
}
func GetWithHeader(ctx context.Context, url string, headers map[string]string, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, sClient, get, url, nil, headers, options...)
}
func GetAndDecode(ctx context.Context, url string, result interface{}, options ...func(context.Context, string, map[string]interface{})) error {
	return GetWithHeaderAndDecode(ctx, url, nil, nil, result, options...)
}
func GetWithHeaderAndDecode(ctx context.Context, url string, obj interface{}, headers map[string]string, result interface{}, options ...func(context.Context, string, map[string]interface{})) error {
	decoder, er1 := DoWithClient(ctx, sClient, get, url, obj, headers, options...)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func Delete(ctx context.Context, url string, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, sClient, delete, url, nil, nil, options...)
}
func DeleteWithHeader(ctx context.Context, url string, headers map[string]string, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, sClient, delete, url, nil, headers, options...)
}
func DeleteAndDecode(ctx context.Context, url string, result interface{}, options ...func(context.Context, string, map[string]interface{})) error {
	return DeleteWithHeaderAndDecode(ctx, url, nil, nil, result, options...)
}
func DeleteWithHeaderAndDecode(ctx context.Context, url string, obj interface{}, headers map[string]string, result interface{}, options ...func(context.Context, string, map[string]interface{})) error {
	decoder, er1 := DoWithClient(ctx, sClient, delete, url, obj, headers, options...)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func Post(ctx context.Context, url string, obj interface{}, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, sClient, post, url, obj, nil, options...)
}
func PostWithHeader(ctx context.Context, url string, obj interface{}, headers map[string]string, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, sClient, post, url, obj, headers, options...)
}
func PostAndDecode(ctx context.Context, url string, obj interface{}, result interface{}, options ...func(context.Context, string, map[string]interface{})) error {
	return PostWithHeaderAndDecode(ctx, url, obj, nil, result, options...)
}
func PostWithHeaderAndDecode(ctx context.Context, url string, obj interface{}, headers map[string]string, result interface{}, options ...func(context.Context, string, map[string]interface{})) error {
	decoder, er1 := DoWithClient(ctx, sClient, post, url, obj, headers, options...)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func Put(ctx context.Context, url string, obj interface{}, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, sClient, put, url, obj, nil, options...)
}
func PutWithHeader(ctx context.Context, url string, obj interface{}, headers map[string]string, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, sClient, put, url, obj, headers, options...)
}
func PutAndDecode(ctx context.Context, url string, obj interface{}, result interface{}, options ...func(context.Context, string, map[string]interface{})) error {
	return PutWithHeaderAndDecode(ctx, url, obj, nil, result, options...)
}
func PutWithHeaderAndDecode(ctx context.Context, url string, obj interface{}, headers map[string]string, result interface{}, options ...func(context.Context, string, map[string]interface{})) error {
	decoder, er1 := DoWithClient(ctx, sClient, put, url, obj, headers, options...)
	if er1 != nil {
		return er1
	}
	er2 := decoder.Decode(result)
	return er2
}
func Patch(ctx context.Context, url string, obj interface{}, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, sClient, patch, url, obj, nil, options...)
}
func PatchWithHeader(ctx context.Context, url string, obj interface{}, headers map[string]string, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	return DoWithClient(ctx, sClient, patch, url, obj, headers, options...)
}
func PatchAndDecode(ctx context.Context, url string, obj interface{}, result interface{}, options ...func(context.Context, string, map[string]interface{})) error {
	return PatchWithHeaderAndDecode(ctx, url, obj, nil, result, options...)
}
func PatchWithHeaderAndDecode(ctx context.Context, url string, obj interface{}, headers map[string]string, result interface{}, options ...func(context.Context, string, map[string]interface{})) error {
	decoder, er1 := DoWithClient(ctx, sClient, patch, url, obj, headers, options...)
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
func DoWithClient(ctx context.Context, client *http.Client, method string, url string, obj interface{}, headers map[string]string, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	rq, err := Marshal(obj)
	if err != nil {
		return nil, err
	}
	return DoAndBuildDecoder(ctx, client, url, method, rq, headers, options...)
}
func DoAndBuildDecoder(ctx context.Context, client *http.Client, url string, method string, body []byte, headers map[string]string, options ...func(context.Context, string, map[string]interface{})) (*json.Decoder, error) {
	var logInfo func(context.Context, string, map[string]interface{})
	if len(options) > 0 {
		logInfo = options[0]
	}
	if conf.Log == true && logInfo != nil {
		if conf.Separate && len(conf.Request) > 0 && body != nil {
			fs1 := make(map[string]interface{}, 0)
			rq := string(body)
			if len(rq) > 0 {
				fs1[conf.Request] = rq
			}
			logInfo(ctx, method+" "+url, fs1)
		}
		start := time.Now()
		res, er1 := DoJSON(ctx, client, url, method, body, headers)
		end := time.Now()
		fs3 := make(map[string]interface{}, 0)
		fs3[conf.Duration] = end.Sub(start).Milliseconds()
		if !conf.Separate && len(conf.Request) > 0 && body != nil {
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
		if len(conf.Size) > 0 {
			fs3[conf.Size] = res.ContentLength
		}
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
		if len(conf.Response) > 0 {
			fs3[conf.Response] = s
		}
		if res.StatusCode == 503 {
			logInfo(ctx, method+" "+url, fs3)
			er2 := errors.New("503 Service Unavailable")
			return nil, er2
		}
		logInfo(ctx, method+" "+url, fs3)
		return json.NewDecoder(strings.NewReader(s)), nil
	} else {
		res, er1 := DoJSON(ctx, client, url, method, body, headers)
		if er1 != nil {
			return nil, er1
		}
		if res.StatusCode == 503 {
			er2 := errors.New("503 Service Unavailable")
			return nil, er2
		}
		return json.NewDecoder(res.Body), nil
	}
}
