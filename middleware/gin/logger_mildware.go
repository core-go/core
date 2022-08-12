package gin

import (
	"context"
	"net"
	"net/http"
	"strings"
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
	if c.Headers != nil && len(c.Headers) > 0 {
		fieldConfig.Headers = c.Headers
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


type ctxKeyRequestID int

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey ctxKeyRequestID = 0

// GetReqID returns a request ID from the given context if one is present.
// Returns the empty string if a request ID cannot be found.
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
