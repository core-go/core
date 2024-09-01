package gin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MaskLogger struct {
	send         func(context.Context, []byte, map[string]string) error
	KeyMap       map[string]string
	RequestKey   string
	MaskRequest  func(map[string]interface{})
	MaskResponse func(map[string]interface{})
	JsonFormat   bool
}

func NewMaskLogger(requestKey string, maskRequest func(map[string]interface{}), maskResponse func(map[string]interface{}), opts ...bool) *MaskLogger {
	jsonFormat := false
	if len(opts) > 0 {
		jsonFormat = opts[0]
	}
	return &MaskLogger{RequestKey: requestKey, MaskRequest: maskRequest, MaskResponse: maskResponse, JsonFormat: jsonFormat}
}
func NewMaskLoggerWithSending(requestKey string, maskRequest func(map[string]interface{}), maskResponse func(map[string]interface{}), jsonFormat bool, send func(context.Context, []byte, map[string]string) error, options ...map[string]string) *MaskLogger {
	var keyMap map[string]string
	if len(options) >= 1 {
		keyMap = options[0]
	}
	return &MaskLogger{RequestKey: requestKey, MaskRequest: maskRequest, MaskResponse: maskResponse, JsonFormat: jsonFormat, send: send, KeyMap: keyMap}
}

func (l *MaskLogger) LogResponse(log func(context.Context, string, map[string]interface{}), r *http.Request, ww ResponseWriter,
	c LogConfig, t1 time.Time, response string, fields map[string]interface{}, includeRequest bool) {
	if includeRequest && len(c.Request) > 0 {
		MaskRequest(c.Request, fields, l.MaskRequest, l.JsonFormat)
	}
	MaskResponse(ww, c, t1, response, fields, l.MaskResponse, l.JsonFormat)
	msg := r.Method + " " + r.RequestURI
	log(r.Context(), msg, fields)
	if l.send != nil {
		go Send(r.Context(), l.send, msg, fields, l.KeyMap)
	}
}
func (l *MaskLogger) LogRequest(log func(context.Context, string, map[string]interface{}), r *http.Request, fields map[string]interface{}) {
	MaskRequest(l.RequestKey, fields, l.MaskRequest, l.JsonFormat)
	msg := "Request " + r.Method + " " + r.RequestURI
	log(r.Context(), msg, fields)
	if l.send != nil {
		go Send(r.Context(), l.send, msg, fields, l.KeyMap)
	}
}

func MaskResponse(ww ResponseWriter, c LogConfig, t1 time.Time, response string, fields map[string]interface{}, mask func(map[string]interface{}), isJsonFormat bool) {
	if len(c.Response) > 0 {
		fields[c.Response] = response
		responseBody := response
		responseMap := map[string]interface{}{}
		json.Unmarshal([]byte(responseBody), &responseMap)
		if len(responseMap) > 0 {
			mask(responseMap)
			if isJsonFormat {
				fields[c.Response] = responseMap
			} else {
				responseString, err := json.Marshal(responseMap)
				if err != nil {
					fmt.Printf("Error: %s", err.Error())
				} else {
					fields[c.Response] = string(responseString)
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
		fields[c.Size] = ww.Size()
	}
}
func MaskRequest(request string, fields map[string]interface{}, mask func(map[string]interface{}), isJsonFormat bool) {
	if len(request) > 0 {
		req, ok := fields[request]
		if ok {
			requestBody, ok2 := req.(string)
			if ok2 {
				requestMap := map[string]interface{}{}
				json.Unmarshal([]byte(requestBody), &requestMap)
				if len(requestMap) > 0 {
					mask(requestMap)
					if isJsonFormat {
						fields[request] = requestMap
					} else {
						requestString, err := json.Marshal(requestMap)
						if err != nil {
							fmt.Printf("Error: %s", err.Error())
						} else {
							fields[request] = string(requestString)
						}
					}
				}
			}
		}
	}
}
