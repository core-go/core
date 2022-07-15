package echo

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func BuildContextWithMask(next http.Handler, mask func(fieldName, s string) string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context
		ctx = r.Context()
		if len(fieldConfig.Ip) > 0 {
			ip := getRemoteIp(r)
			ctx = context.WithValue(ctx, fieldConfig.Ip, ip)
		}
		if fieldConfig.Constants != nil && len(fieldConfig.Constants) > 0 {
			for k, e := range fieldConfig.Constants {
				if len(e) > 0 {
					ctx = context.WithValue(ctx, k, e)
				}
			}
		}
		if fieldConfig.Map != nil && len(fieldConfig.Map) > 0 && r.Body != nil && (r.Method != "GET" || r.Method != "DELETE") {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			r.Body = ioutil.NopCloser(buf)
			var v interface{}
			er2 := json.NewDecoder(strings.NewReader(buf.String())).Decode(&v)
			if er2 != nil {
				if len(fieldConfig.Ip) == 0 && fieldConfig.Constants == nil {
					next.ServeHTTP(w, r)
				} else {
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			} else {
				m, ok := v.(map[string]interface{})
				if !ok {
					if len(fieldConfig.Ip) == 0 && fieldConfig.Constants == nil {
						next.ServeHTTP(w, r)
					} else {
						next.ServeHTTP(w, r.WithContext(ctx))
					}
				} else {
					for k, e := range fieldConfig.Map {
						if strings.Index(e, ".") >= 0 {
							v3 := ValueOf(v, e)
							if v3 != nil {
								s3, ok3 := v3.(string)
								if ok3 {
									if len(s3) > 0 {
										if mask != nil && fieldConfig.Masks != nil && len(fieldConfig.Masks) > 0 {
											if Include(fieldConfig.Masks, k) {
												ctx = context.WithValue(ctx, k, mask(k, s3))
											} else {
												ctx = context.WithValue(ctx, k, s3)
											}
										} else {
											ctx = context.WithValue(ctx, k, s3)
										}
									}
								} else {
									ctx = context.WithValue(ctx, k, v3)
								}
							}
						} else {
							x, ok2 := m[e]
							if ok2 && x != nil {
								s3, ok3 := x.(string)
								if ok3 {
									if len(s3) > 0 {
										if mask != nil && fieldConfig.Masks != nil && len(fieldConfig.Masks) > 0 {
											if Include(fieldConfig.Masks, k) {
												ctx = context.WithValue(ctx, k, mask(k, s3))
											} else {
												ctx = context.WithValue(ctx, k, s3)
											}
										} else {
											ctx = context.WithValue(ctx, k, s3)
										}
									}
								} else {
									ctx = context.WithValue(ctx, k, x)
								}
							}
						}
					}
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		} else {
			if len(fieldConfig.Ip) == 0 && fieldConfig.Constants == nil {
				next.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		}
	}
	return http.HandlerFunc(fn)
}
func BuildContext(next http.Handler) http.Handler {
	return BuildContextWithMask(next, nil)
}
func Include(vs []string, v string) bool {
	for _, s := range vs {
		if v == s {
			return true
		}
	}
	return false
}
func ValueOf(m interface{}, path string) interface{} {
	arr := strings.Split(path, ".")
	i := 0
	var c interface{}
	c = m
	l1 := len(arr) - 1
	for i < len(arr) {
		key := arr[i]
		m2, ok := c.(map[string]interface{})
		if ok {
			c = m2[key]
		}
		if !ok || i >= l1 {
			return c
		}
		i++
	}
	return c
}
