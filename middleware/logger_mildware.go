package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"
)

func InitializeFieldConfig(c LogConfig) {
	if len(c.Duration) > 0 {
		fieldConfig.Duration = c.Duration
	} else {
		fieldConfig.Duration = "duration"
	}
	fieldConfig.Log = c.Log
	fieldConfig.Ip = c.Ip
	if c.Map != nil && len(c.Map) > 0 {
		fieldConfig.Map = c.Map
	}
	if c.Constants != nil && len(c.Constants) > 0 {
		fieldConfig.Constants = c.Constants
	}
	if len(c.Fields) > 0 {
		fields := strings.Split(c.Fields, ",")
		fieldConfig.Fields = fields
	}
	if len(c.Masks) > 0 {
		fields := strings.Split(c.Masks, ",")
		fieldConfig.Masks = fields
	}
	if len(c.Skips) > 0 {
		fields := strings.Split(c.Skips, ",")
		fieldConfig.Skips = fields
	}
}
func Logger(c LogConfig, log func(ctx context.Context, msg string, fields map[string]interface{}), f Formatter) func(h http.Handler) http.Handler {
	InitializeFieldConfig(c)
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if !fieldConfig.Log || InSkipList(r, fieldConfig.Skips) {
				h.ServeHTTP(w, r)
			} else {
				dw := NewResponseWriter(w)
				ww := NewWrapResponseWriter(dw, r.ProtoMajor)
				startTime := time.Now()
				fields := BuildLogFields(c, r)
				single := !c.Separate
				if r.Method == "GET" || r.Method == "DELETE" {
					single = true
				}
				f.LogRequest(log, r, c, fields, single)
				defer func() {
					if single {
						f.LogResponse(log, r, ww, c, startTime, dw.Body.String(), fields, single)
					} else {
						resLogFields := BuildLogFields(c, r)
						f.LogResponse(log, r, ww, c, startTime, dw.Body.String(), resLogFields, single)
					}
				}()
				h.ServeHTTP(ww, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}
func InSkipList(r *http.Request, skips []string) bool {
	if skips == nil || len(skips) == 0 {
		return false
	}
	for _, s := range skips {
		if strings.HasSuffix(s, r.RequestURI) {
			return true
		}
	}
	return false
}
func BuildLogFields(c LogConfig, r *http.Request) map[string]interface{} {
	fields := make(map[string]interface{}, 0)
	if !c.Build {
		return fields
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if len(c.Uri) > 0 {
		fields[c.Uri] = r.RequestURI
	}

	if len(c.ReqId) > 0 {
		if reqID := GetReqID(r.Context()); reqID != "" {
			fields[c.ReqId] = reqID
		}
	}
	if len(c.Scheme) > 0 {
		fields[c.Scheme] = scheme
	}
	if len(c.Proto) > 0 {
		fields[c.Proto] = r.Proto
	}
	if len(c.UserAgent) > 0 {
		fields[c.UserAgent] = r.UserAgent()
	}
	if len(c.RemoteAddr) > 0 {
		fields[c.RemoteAddr] = r.RemoteAddr
	}
	if len(c.Method) > 0 {
		fields[c.Method] = r.Method
	}
	if len(c.RemoteIp) > 0 {
		remoteIP := getRemoteIp(r)
		fields[c.RemoteIp] = remoteIP
	}
	return fields
}
func getRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
