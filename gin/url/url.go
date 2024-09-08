package url

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
func QueryRequiredString(c *gin.Context, name string) string {
	s := QueryString(c.Request.URL.Query(), name)
	if len(s) == 0 {
		c.String(http.StatusBadRequest, fmt.Sprintf("%s is required", name))
	}
	return s
}
func QueryRequiredStrings(c *gin.Context, name string, opts ...string) []string {
	s := QueryString(c.Request.URL.Query(), name)
	if len(s) == 0 {
		c.String(http.StatusBadRequest, fmt.Sprintf("%s is required", name))
		return nil
	} else {
		if len(opts) > 0 && len(opts[0]) > 0 {
			return strings.Split(s, opts[0])
		} else {
			return strings.Split(s, ",")
		}
	}
}
func QueryRequiredInt64(c *gin.Context, name string) (int64, error) {
	s := QueryString(c.Request.URL.Query(), name)
	if len(s) > 0 {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "parameter must be an integer")
			return 0, err
		}
		return i, nil
	}
	se := fmt.Sprintf("%s is a required integer", name)
	c.String(http.StatusBadRequest, se)
	return 0, errors.New(se)
}
func QueryRequiredUint64(c *gin.Context, name string) (uint64, error) {
	s := QueryString(c.Request.URL.Query(), name)
	if len(s) > 0 {
		i, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "parameter must be an unsigned integer")
			return 0, err
		}
		return i, nil
	}
	se := fmt.Sprintf("%s is a required integer", name)
	c.String(http.StatusBadRequest, se)
	return 0, errors.New(se)
}
func QueryRequiredInt32(c *gin.Context, name string) (int32, error) {
	s := QueryString(c.Request.URL.Query(), name)
	if len(s) > 0 {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "parameter must be an integer")
			return 0, err
		}
		return int32(i), nil
	}
	se := fmt.Sprintf("%s is a required integer", name)
	c.String(http.StatusBadRequest, se)
	return 0, errors.New(se)
}
func QueryRequiredInt(c *gin.Context, name string) (int, error) {
	s := QueryString(c.Request.URL.Query(), name)
	if len(s) > 0 {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "parameter must be an integer")
			return 0, err
		}
		return int(i), nil
	}
	se := fmt.Sprintf("%s is a required integer", name)
	c.String(http.StatusBadRequest, se)
	return 0, errors.New(se)
}
