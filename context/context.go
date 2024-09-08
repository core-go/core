package context

import "context"

func FromContext(ctx context.Context, key string, options ...string) string {
	var authorization string
	if len(options) > 0 {
		authorization = options[0]
	}
	if len(authorization) > 0 {
		token := ctx.Value(authorization)
		if token != nil {
			if authorizationToken, exist := token.(map[string]interface{}); exist {
				return FromMap(key, authorizationToken)
			}
		}
		return ""
	} else {
		u := ctx.Value(key)
		if u != nil {
			v, ok := u.(string)
			if ok {
				return v
			}
		}
		return ""
	}
}
func FromMap(key string, data map[string]interface{}) string {
	if data == nil {
		return ""
	}
	u := data[key]
	if u != nil {
		v, ok := u.(string)
		if ok {
			return v
		}
	}
	return ""
}
