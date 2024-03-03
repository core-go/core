package gin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type MaskLogger struct {
	send       func(ctx context.Context, data []byte, attributes map[string]string) (string, error)
	KeyMap     map[string]string
	Goroutines bool
	MaskRequest func(fieldName string, v interface{}) interface{}
	MaskResponse func(fieldName string, v interface{}) interface{}
}
func NewMaskLogger(maskRequest func(fieldName string, v interface{}) interface{}, maskResponse func(fieldName string, v interface{}) interface{}) *MaskLogger {
	return &MaskLogger{MaskRequest: maskRequest, MaskResponse: maskResponse}
}
func NewMaskLoggerWithSending(maskRequest func(fieldName string, v interface{}) interface{}, maskResponse func(fieldName string, v interface{}) interface{}, send func(context.Context, []byte, map[string]string) (string, error), goroutines bool, options ...map[string]string) *MaskLogger {
	var keyMap map[string]string
	if len(options) >= 1 {
		keyMap = options[0]
	}
	return &MaskLogger{MaskRequest: maskRequest, MaskResponse: maskResponse, send: send, Goroutines: goroutines, KeyMap: keyMap}
}

func (l *MaskLogger) LogResponse(log func(context.Context, string, map[string]interface{}), r *http.Request, ww ResponseWriter,
	c LogConfig, t1 time.Time, response string, fields map[string]interface{}, singleLog bool) {
	fs := BuildMaskedResponseBody(ww, c, t1, response, fields, l.MaskResponse)
	var msg string
	if singleLog {
		msg = r.Method + " " + r.RequestURI
	} else {
		msg = "Response " + r.Method + " " + r.RequestURI
	}
	log(r.Context(), msg, fs)
	if l.send != nil {
		if l.Goroutines {
			go Send(r.Context(), l.send, msg, fields, l.KeyMap)
		} else {
			Send(r.Context(), l.send, msg, fields, l.KeyMap)
		}
	}
}
func (l *MaskLogger) LogRequest(log func(context.Context, string, map[string]interface{}), r *http.Request, c LogConfig, fields map[string]interface{}, singleLog bool) {
	var fs map[string]interface{}
	fs = fields
	if len(c.Request) > 0 && r.Method != "GET" && r.Method != "DELETE" && !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		fs = BuildMaskedRequestBody(r, c.Request, fields, l.MaskRequest)
	}
	if !singleLog {
		msg := "Request " + r.Method + " " + r.RequestURI
		log(r.Context(), msg, fs)
		if l.send != nil {
			if l.Goroutines {
				go Send(r.Context(), l.send, msg, fields, l.KeyMap)
			} else {
				Send(r.Context(), l.send, msg, fields, l.KeyMap)
			}
		}
	}
}

func BuildMaskedResponseBody(ww ResponseWriter, c LogConfig, t1 time.Time, response string, fields map[string]interface{}, mask func(fieldName string, s interface{}) interface{}) map[string]interface{} {
	if len(c.Response) > 0 {
		fields[c.Response] = response
		responseBody := response
		responseMap := map[string]interface{}{}
		json.Unmarshal([]byte(responseBody), &responseMap)
		if len(responseMap) > 0 {
			for key, v := range responseMap {
				responseMap[key] = mask(key, v)
			}
			responseString, err :=  json.Marshal(responseMap)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
			} else {
				fields[c.Response] = string(responseString)
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
		fields[c.Size] = ww.Size()
	}
	return fields
}
func BuildMaskedRequestBody(r *http.Request, request string, fields map[string]interface{}, mask func(fieldName string, s interface{}) interface{}) map[string]interface{} {
	if r.Body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		fields[request] = buf.String()
		r.Body = ioutil.NopCloser(buf)
		requestBody := fields[request].(string)
		requestMap := map[string]interface{}{}
		json.Unmarshal([]byte(requestBody), &requestMap)
		if len(requestMap) > 0 {
			for key, v := range requestMap {
				requestMap[key] = mask(key, v)
			}
			requestString, err :=  json.Marshal(requestMap)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
			} else {
				fields[request] = string(requestString)
			}
		}
	}
	return fields
}
