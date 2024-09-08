package url

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func QueryString(v url.Values, name string, opts ...string) string {
	s := v.Get(name)
	if len(s) > 0 {
		return s
	}
	if len(opts) > 0 {
		return opts[0]
	}
	return ""
}
func QueryStrings(v url.Values, name string, opts ...[]string) []string {
	s, ok := v[name]
	if ok {
		return s
	}
	if len(opts) > 0 {
		return opts[0]
	}
	return nil
}
func QueryArray(v url.Values, name string, all []string, opts ...[]string) []string {
	s, ok := v[name]
	if ok {
		x := QueryIn(all, s)
		return x
	}
	if len(opts) > 0 {
		return opts[0]
	}
	return nil
}
func isIn(arr []string, s string) bool {
	for _, a := range arr {
		if s == a {
			return true
		}
	}
	return false
}
func QueryIn(all []string, s []string) []string {
	var fieldsParamArr []string
	checkSubstr := strings.Index(s[0], ",")
	if checkSubstr > 0 {
		fieldsParamArr = strings.Split(s[0], ",")
	} else {
		fieldsParamArr = s
	}
	for _, v := range fieldsParamArr {
		valueTrim := strings.TrimSpace(v)
		check := isIn(all, valueTrim)
		if check == false {
			return nil
		}
	}
	return fieldsParamArr
}

func QueryInt64(v url.Values, name string, opts ...int64) *int64 {
	s := QueryString(v, name)
	if len(s) > 0 {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil
		}
		return &i
	}
	if len(opts) > 0 {
		return &opts[0]
	}
	return nil
}
func QueryInt32(v url.Values, name string, opts ...int64) *int32 {
	i := QueryInt64(v, name, opts...)
	if i != nil {
		j := int32(*i)
		return &j
	}
	return nil
}
func QueryInt(v url.Values, name string, opts ...int64) *int {
	i := QueryInt64(v, name, opts...)
	if i != nil {
		j := int(*i)
		return &j
	}
	return nil
}
func QueryRequiredString(w http.ResponseWriter, v url.Values, name string) string {
	s := QueryString(v, name)
	if len(s) == 0 {
		http.Error(w, fmt.Sprintf("%s is required", name), http.StatusBadRequest)
	}
	return s
}
func QueryRequiredStrings(w http.ResponseWriter, v url.Values, name string, opts ...string) []string {
	s := QueryString(v, name)
	if len(s) == 0 {
		http.Error(w, fmt.Sprintf("%s is required", name), http.StatusBadRequest)
		return nil
	} else {
		if len(opts) > 0 && len(opts[0]) > 0 {
			return strings.Split(s, opts[0])
		} else {
			return strings.Split(s, ",")
		}
	}
}
func QueryRequiredInt64(w http.ResponseWriter, v url.Values, name string) (int64, error) {
	s := QueryString(v, name)
	if len(s) > 0 {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 0, err
		}
		return i, nil
	}
	se := fmt.Sprintf("%s is a required integer", name)
	http.Error(w, se, http.StatusBadRequest)
	return 0, errors.New(se)
}
func QueryRequiredInt32(w http.ResponseWriter, v url.Values, name string) (int32, error) {
	s := QueryString(v, name)
	if len(s) > 0 {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 0, err
		}
		return int32(i), nil
	}
	se := fmt.Sprintf("%s is a required integer", name)
	http.Error(w, se, http.StatusBadRequest)
	return 0, errors.New(se)

}
func QueryRequiredInt(w http.ResponseWriter, v url.Values, name string) (int, error) {
	s := QueryString(v, name)
	if len(s) > 0 {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 0, err
		}
		return int(i), nil
	}
	se := fmt.Sprintf("%s is a required integer", name)
	http.Error(w, se, http.StatusBadRequest)
	return 0, errors.New(se)
}
