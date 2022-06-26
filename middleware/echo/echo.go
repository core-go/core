package echo

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/core-go/core/middleware"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"strings"
	"time"
)

type Fields map[string]interface{}

type Logger struct {
	Config  middleware.LogConfig
	LogInfo func(ctx context.Context, msg string, fields map[string]interface{})
	f       middleware.Formatter
	Mask    func(fieldName, s string) string
}

func NewLogger(c middleware.LogConfig, log func(ctx context.Context, msg string, fields map[string]interface{}), f middleware.Formatter, mask func(fieldName, s string) string) *Logger {
	return &Logger{c, log, f, mask}
}

var fieldConfig middleware.FieldConfig

func InitializeFieldConfig(c middleware.LogConfig) {
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
func (l *Logger) LoggerEcho(next echo.HandlerFunc) echo.HandlerFunc {
	InitializeFieldConfig(l.Config)
	return func(c echo.Context) error {
		if !fieldConfig.Log || middleware.InSkipList(c.Request(), fieldConfig.Skips) {
			return next(c)
		} else {
			r := c.Request()
			dw := middleware.NewResponseWriter(c.Response().Writer)
			ww := middleware.NewWrapResponseWriter(dw, r.ProtoMajor)
			startTime := time.Now()
			fields := middleware.BuildLogFields(l.Config, r)
			single := !l.Config.Separate
			if r.Method == "GET" || r.Method == "DELETE" {
				single = true
			}

			l.f.LogRequest(l.LogInfo, r, l.Config, fields, single)
			c.Response().Writer = ww
			defer func() {
				if single {
					l.f.LogResponse(l.LogInfo, r, ww, l.Config, startTime, dw.Body.String(), fields, single)
				} else {
					resLogFields := middleware.BuildLogFields(l.Config, r)
					l.f.LogResponse(l.LogInfo, r, ww, l.Config, startTime, dw.Body.String(), resLogFields, single)
				}
			}()
			return next(c)
		}
		return nil
	}
}

func (l *Logger) BuildContextWithMask(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctxEcho := c
		var ctx context.Context
		ctx = c.Request().Context()
		if fieldConfig.Constants != nil && len(fieldConfig.Constants) > 0 {
			for k, e := range fieldConfig.Constants {
				if len(e) > 0 {
					ctx = context.WithValue(ctx, k, e)
				}
			}
		}

		r := c.Request()
		if fieldConfig.Map != nil && len(fieldConfig.Map) > 0 && r.Body != nil && (r.Method != "GET" || r.Method != "DELETE") {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			r.Body = ioutil.NopCloser(buf)
			var v interface{}
			er2 := json.NewDecoder(strings.NewReader(buf.String())).Decode(&v)
			if er2 != nil {
				if len(fieldConfig.Ip) == 0 && fieldConfig.Constants == nil {
					next(c)
				} else {
					ctxEcho.Request().WithContext(ctx)
					next(ctxEcho)
				}
			} else {
				m, ok := v.(map[string]interface{})
				if !ok {
					if len(fieldConfig.Ip) == 0 && fieldConfig.Constants == nil {
						return next(c)
					} else {
						ctxEcho.Request().WithContext(ctx)
						return next(ctxEcho)
					}
				} else {
					for k, e := range fieldConfig.Map {
						if strings.Index(e, ".") >= 0 {
							v3 := middleware.ValueOf(v, e)
							if v3 != nil {
								s3, ok3 := v3.(string)
								if ok3 {
									if len(s3) > 0 {
										if l.Mask != nil && fieldConfig.Masks != nil && len(fieldConfig.Masks) > 0 {
											if middleware.Include(fieldConfig.Masks, k) {
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
											if middleware.Include(fieldConfig.Masks, k) {
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
					ctxEcho.SetRequest(c.Request().WithContext(ctx))
					return next(ctxEcho)
				}
			}
		} else {
			if len(fieldConfig.Ip) == 0 && fieldConfig.Constants == nil {
				return next(c)
			} else {
				ctxEcho.SetRequest(c.Request().WithContext(ctx))
				return next(ctxEcho)
			}
		}
		return nil
	}
}
