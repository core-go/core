package gin

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
	"time"
)

type Fields map[string]interface{}

type GinLogger struct {
	Config  LogConfig
	LogInfo func(ctx context.Context, msg string, fields map[string]interface{})
	f       Formatter
	Mask    func(fieldName, s string) string
}

func NewGinLogger(c LogConfig, logInfo func(ctx context.Context, msg string, fields map[string]interface{}), mask func(fieldName, s string) string) *GinLogger {
	logger := NewLogger()
	return &GinLogger{c, logInfo, logger, mask}
}

func (l *GinLogger) Logger() gin.HandlerFunc {
	InitializeFieldConfig(l.Config)
	return func(c *gin.Context) {
		if !fieldConfig.Log || InSkipList(c.Request, fieldConfig.Skips) {
			c.Next()
		} else {
			r := c.Request
			dw := NewResponseWriter(c.Writer)

			startTime := time.Now()
			fields := BuildLogFields(l.Config, r)
			single := !l.Config.Separate
			if r.Method == "GET" || r.Method == "DELETE" {
				single = true
			}

			l.f.LogRequest(l.LogInfo, r, l.Config, fields, single)
			c.Writer = dw
			defer func() {
				if single {
					l.f.LogResponse(l.LogInfo, r, *dw, l.Config, startTime, dw.Body.String(), fields, single)
				} else {
					resLogFields := BuildLogFields(l.Config, r)
					l.f.LogResponse(l.LogInfo, r, *dw, l.Config, startTime, dw.Body.String(), resLogFields, single)
				}
			}()
			c.Next()
		}
	}
}

func (l *GinLogger) BuildContextWithMask() gin.HandlerFunc {
	return func(c *gin.Context)  {
		ctxGin := c
		var ctx context.Context
		ctx = c.Request.Context()
		if fieldConfig.Constants != nil && len(fieldConfig.Constants) > 0 {
			for k, e := range fieldConfig.Constants {
				if len(e) > 0 {
					ctx = context.WithValue(ctx, k, e)
				}
			}
		}

		r := c.Request
		if fieldConfig.Map != nil && len(fieldConfig.Map) > 0 && r.Body != nil && (r.Method != "GET" || r.Method != "DELETE") {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			r.Body = ioutil.NopCloser(buf)
			var v interface{}
			er2 := json.NewDecoder(strings.NewReader(buf.String())).Decode(&v)
			if er2 != nil {
				if len(fieldConfig.Ip) == 0 && fieldConfig.Constants == nil {
					c.Next()
				} else {

					ctxGin.Request.WithContext(ctx)
					ctxGin.Next()
				}
			} else {
				m, ok := v.(map[string]interface{})
				if !ok {
					if len(fieldConfig.Ip) == 0 && fieldConfig.Constants == nil {
						c.Next()
					} else {
						ctxGin.Request.WithContext(ctx)
						ctxGin.Next()
					}
				} else {
					for k, e := range fieldConfig.Map {
						if strings.Index(e, ".") >= 0 {
							v3 := ValueOf(v, e)
							if v3 != nil {
								s3, ok3 := v3.(string)
								if ok3 {
									if len(s3) > 0 {
										if l.Mask != nil && fieldConfig.Masks != nil && len(fieldConfig.Masks) > 0 {
											if Include(fieldConfig.Masks, k) {
												ctx = context.WithValue(ctx, k, l.Mask(k, s3))
											} else {
												ctx = context.WithValue(ctx, k, s3)
											}
										} else {
											ctx = context.WithValue(ctx, k, s3)
										}
									}
								} else {
									ctx = context.WithValue(ctx, k, s3)
								}
							}
						} else {
							x, ok2 := m[e]
							if ok2 && x != nil {
								s3, ok3 := x.(string)
								if ok3 {
									if len(s3) > 0 {
										if l.Mask != nil && fieldConfig.Masks != nil && len(fieldConfig.Masks) > 0 {
											if Include(fieldConfig.Masks, k) {
												ctx = context.WithValue(ctx, k, l.Mask(k, s3))
											} else {
												ctx = context.WithValue(ctx, k, s3)
											}
										} else {
											ctx = context.WithValue(ctx, k, s3)
										}
									}
								} else {
									ctx = context.WithValue(ctx, k, s3)
								}
							}
						}
					}
					ctxGin.Request = ctxGin.Request.WithContext(ctx)
					ctxGin.Next()
				}
			}
		} else {
			if len(fieldConfig.Ip) == 0 && fieldConfig.Constants == nil {
				c.Next()
			} else {
				ctxGin.Request = ctxGin.Request.WithContext(ctx)
				ctxGin.Next()
			}
		}
	}
}
