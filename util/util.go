package util

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/core-go/core/template"
)

const dateFormat = "2006-01-02"

func CreateFuncMap() template.FuncMap {
	funcMap := make(template.FuncMap, 0)
	funcMap["SKIP"] = Skip
	funcMap["ToYesNo"] = ToYesNo
	funcMap["PtrToYesNo"] = PtrToYesNo
	funcMap["ToBool"] = ToBool
	funcMap["PtrToBool"] = PtrToBool
	funcMap["ToString"] = ToString
	funcMap["ParseFloat"] = ParseFloat
	funcMap["ToFloatPtr"] = ToFloatPtr
	funcMap["CheckCase"] = CheckCase
	funcMap["CheckStrings"] = CheckStrings
	funcMap["FormatTime"] = FormatTime
	return funcMap
}
func CreateParam(name string, v interface{}, opts ...map[string]interface{}) map[string]interface{} {
	param := make(map[string]interface{})
	param[name] = v

	t := time.Now()
	date := t.Format(dateFormat)
	param["now"] = t
	param["today"] = date
	if len(opts) > 0 && opts[0] != nil {
		for key, element := range opts[0] {
			param[key] = element
		}
	}
	return param
}
func Unmarshal(t *template.Template, param interface{}, op interface{}) error {
	buf := &bytes.Buffer{}
	err := t.Execute(buf, param)
	if err != nil {
		return err
	}
	return json.NewDecoder(strings.NewReader(buf.String())).Decode(op)
}
func PtrToBool(s *string, v string) bool {
	if s != nil && *s == v {
		return true
	}
	return false
}
func ToBool(s string, v string) bool {
	if s == v {
		return true
	}
	return false
}
func PtrToYesNo(s *string, v string) string {
	if s != nil && *s == v {
		return "Y"
	}
	return "N"
}
func ToYesNo(s string, v string) string {
	if s == v {
		return "Y"
	}
	return "N"
}
func Equal(s *string, v string) bool {
	if s != nil && *s == v {
		return true
	}
	return false
}
func ToString(s *string, v string) string {
	if s != nil && *s != v {
		return *s
	}
	return ""
}
func ParseFloat(s *string, v string) float64 {
	if s != nil && *s != v {
		min, _ := strconv.ParseFloat(*s, 64)
		return min
	}
	return 0
}
func ToFloatPtr(s *string, v string) *float64 {
	if s != nil && *s != v {
		min, _ := strconv.ParseFloat(*s, 64)
		return &min
	}
	return nil
}
func CheckCase(i, c, t, e string) string {
	if i == c {
		return t
	} else {
		return e
	}
}
func CheckStrings(s []string, v string, r int) *int {
	if len(s) == 0 {
		return nil
	}
	count := 0
	for _, v := range s {
		if v != v {
			count++
		}
	}
	if count >= 1 {
		return &r
	}
	return nil
}
func Skip(v interface{}) interface{} {
	return v
}
func FormatTime(timeValue time.Time) string {
	return timeValue.Format(time.RFC3339)
}
