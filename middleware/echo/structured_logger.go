package echo

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Formatter interface {
	LogRequest(log func(context.Context, string, map[string]interface{}), r *http.Request, fields map[string]interface{})
	LogResponse(log func(context.Context, string, map[string]interface{}), r *http.Request, ww WrapResponseWriter, c LogConfig, startTime time.Time, response string, fields map[string]interface{}, includeRequest bool)
}
type StructuredLogger struct {
	send       func(context.Context, []byte, map[string]string) error
	KeyMap     map[string]string
	RequestKey string
	JsonFormat bool
}

var fieldConfig FieldConfig

func NewLogger() *StructuredLogger {
	return &StructuredLogger{}
}
func NewLoggerWithJsonFormat(requestKey string, jsonFormat bool) *StructuredLogger {
	return &StructuredLogger{RequestKey: requestKey, JsonFormat: jsonFormat}
}
func NewLoggerWithSending(requestKey string, jsonFormat bool, send func(context.Context, []byte, map[string]string) error, options ...map[string]string) *StructuredLogger {
	var keyMap map[string]string
	if len(options) >= 1 {
		keyMap = options[0]
	}
	return &StructuredLogger{RequestKey: requestKey, JsonFormat: jsonFormat, send: send, KeyMap: keyMap}
}
func (l *StructuredLogger) LogResponse(log func(context.Context, string, map[string]interface{}), r *http.Request, ww WrapResponseWriter,
	c LogConfig, t1 time.Time, response string, fields map[string]interface{}, includeRequest bool) {
	BuildResponse(ww, c, t1, response, fields, l.JsonFormat)
	msg := r.Method + " " + r.RequestURI
	log(r.Context(), msg, fields)
	if l.send != nil {
		go Send(r.Context(), l.send, msg, fields, l.KeyMap)
	}
}
func Send(ctx context.Context, send func(context.Context, []byte, map[string]string) error, msg string, fields map[string]interface{}, keyMap map[string]string) {
	m2 := AddKeyFields(msg, fields, keyMap)
	b, err := json.Marshal(m2)
	if err == nil {
		send(ctx, b, nil)
	}
}
func (l *StructuredLogger) LogRequest(log func(context.Context, string, map[string]interface{}), r *http.Request, fields map[string]interface{}) {
	msg := "Request " + r.Method + " " + r.RequestURI
	if l.JsonFormat && len(l.RequestKey) > 0 {
		req, ok := fields[l.RequestKey]
		if ok {
			requestBody, ok2 := req.(string)
			if ok2 {
				requestMap := map[string]interface{}{}
				json.Unmarshal([]byte(requestBody), &requestMap)
				if len(requestMap) > 0 {
					fields[l.RequestKey] = requestMap
				}
			}
		}
	}
	log(r.Context(), msg, fields)
	if l.send != nil {
		go Send(r.Context(), l.send, msg, fields, l.KeyMap)
	}
}

func BuildResponse(ww WrapResponseWriter, c LogConfig, t1 time.Time, response string, fields map[string]interface{}, jsonFormat bool) {
	if len(c.Response) > 0 {
		if jsonFormat {
			responseBody := response
			responseMap := map[string]interface{}{}
			json.Unmarshal([]byte(responseBody), &responseMap)
			if len(responseMap) > 0 {
				fields[c.Response] = responseMap
			} else {
				fields[c.Response] = response
			}
		} else {
			fields[c.Response] = response
		}
	}
	if jsonFormat && len(c.Request) > 0 {
		req, ok := fields[c.Request]
		if ok {
			requestBody, ok2 := req.(string)
			if ok2 {
				requestMap := map[string]interface{}{}
				json.Unmarshal([]byte(requestBody), &requestMap)
				if len(requestMap) > 0 {
					fields[c.Request] = requestMap
				}
			}
		}
	}
	if len(c.ResponseStatus) > 0 {
		fields[c.ResponseStatus] = ww.Status()
	}
	if len(fieldConfig.Duration) > 0 {
		t2 := time.Now()
		duration := t2.Sub(t1)
		fields[fieldConfig.Duration] = duration.Milliseconds()
	}
	if len(c.Size) > 0 {
		fields[c.Size] = ww.BytesWritten()
	}
}
func BuildRequest(r *http.Request, request string, fields map[string]interface{}) map[string]interface{} {
	if r.Body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		fields[request] = buf.String()
		r.Body = io.NopCloser(buf)
	}
	return fields
}
func AddKeyFields(message string, m map[string]interface{}, keys map[string]string) map[string]interface{} {
	level := "level"
	t := "time"
	msg := "msg"
	if keys != nil {
		ks := keys
		v1, ok1 := ks[level]
		if ok1 && len(v1) > 0 {
			level = v1
		}
		v2, ok2 := ks[t]
		if ok2 && len(v2) > 0 {
			t = v2
		}
		v3, ok3 := ks[msg]
		if ok3 && len(v3) > 0 {
			msg = v3
		}
	}
	m[msg] = message
	m[level] = "info"
	m[t] = time.Now()
	return m
}
