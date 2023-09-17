package gin

import (
	"fmt"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func GetPositionFromSortedPrivileges(privileges []string, privilegeId string) int {
	return sort.Search(len(privileges), func(i int) bool { return privileges[i][:strings.Index(privileges[i], ":")] >= privilegeId })
}

func GetAction(privileges []string, privilegeId string, sortedPrivilege bool) int32 {
	prefix := fmt.Sprintf("%s:", privilegeId)
	prefixLen := len(prefix)

	if sortedPrivilege {
		i := GetPositionFromSortedPrivileges(privileges, privilegeId)
		if i >= 0 && i < len(privileges) && strings.HasPrefix(privileges[i], prefix) {
			hexAction := privileges[i][prefixLen:]
			if action, err := ConvertHexAction(hexAction); err == nil {
				return action
			}
		}
	} else {
		for _, privilege := range privileges {
			if strings.HasPrefix(privilege, prefix) {
				hexAction := privilege[prefixLen:]
				if action, err := ConvertHexAction(hexAction); err == nil {
					return action
				}
			}
		}
	}

	return ActionNone
}

func ConvertHexAction(hexAction string) (int32, error) {
	if action64, err := strconv.ParseInt(hexAction, 16, 64); err == nil {
		return int32(action64), nil
	} else {
		return 0, err
	}
}
func Include(vs []string, v string) bool {
	for _, s := range vs {
		if v == s {
			return true
		}
	}
	return false
}
func IncludeOfSort(vs []string, v string) bool {
	i := sort.SearchStrings(vs, v)
	if i >= 0 && vs[i] == v {
		return true
	}
	return false
}
func ValueFromMap(key string, data map[string]interface{}) string {
	u := data[key]
	if u != nil {
		v, ok := u.(string)
		if ok {
			return v
		}
	}
	return ""
}
func FromContext(r *http.Request, authorization string, key string) string {
	if len(authorization) > 0 {
		token := r.Context().Value(authorization)
		if token != nil {
			if authorizationToken, exist := token.(map[string]interface{}); exist {
				return ValueFromMap(key, authorizationToken)
			}
		}
		return ""
	} else {
		u := r.Context().Value(key)
		if u != nil {
			v, ok := u.(string)
			if ok {
				return v
			}
		}
		return ""
	}
}
func ValuesFromContext(r *http.Request, authorization string, key string) *[]string {
	token := r.Context().Value(authorization)
	if token != nil {
		if authorizationToken, exist := token.(map[string]interface{}); exist {
			v, ok2 := authorizationToken[key]
			if !ok2 || v == nil {
				return nil
			}
			if values, ok3 := v.(*[]string); ok3 {
				return values
			}
		}
	}
	return nil
}

func getRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}
