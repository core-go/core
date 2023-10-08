package validator

import (
	s "github.com/core-go/core"
	"github.com/go-playground/validator/v10"
	"time"
)

type CustomValidate struct {
	Fn  validator.Func
	Tag string
}

var PatternMap = map[string]string{
	"digit":      "^\\d+$",
	"dash_digit": "^[0-9-]*$",
	"code":       "^\\w*\\d*$",
}

var translations = map[string]string{
	"email":        "{0} must be a valid email address",
	"url":          "{0} must be a valid URL",
	"uri":          "{0} must be a valid URI",
	"fax":          "{0} must be a valid fax number",
	"phone":        "{0} must be a valid phone number",
	"ip":           "{0} must be a valid IP address",
	"ipv4":         "{0} must be a valid IPv4 address",
	"ipv6":         "{0} must be a valid IPv6 address",
	"digit":        "{0} must contain only digits",
	"pin":          "{0} must contain only digits",
	"abc":          "{0} must contain only letters",
	"id":           "{0} must be a valid ID",
	"code":         "{0} must be a valid code",
	"country_code": "{0} must be a valid country code",
	"username":     "{0} must be a valid username",
	"regex":        "{0} must match the provided regex pattern",
	"after_now":    "{0} must be after now",
	"now_or_after": "{0} must be now or after",
}

func GetCustomValidateList() (list []CustomValidate) {
	list = append(list, CustomValidate{Fn: CheckEmail, Tag: "email"})
	list = append(list, CustomValidate{Fn: CheckUrl, Tag: "url"})
	list = append(list, CustomValidate{Fn: CheckUri, Tag: "uri"})
	list = append(list, CustomValidate{Fn: CheckFax, Tag: "fax"})
	list = append(list, CustomValidate{Fn: CheckPhone, Tag: "phone"})
	list = append(list, CustomValidate{Fn: CheckIp, Tag: "ip"})
	list = append(list, CustomValidate{Fn: CheckIpV4, Tag: "ipv4"})
	list = append(list, CustomValidate{Fn: CheckIpV6, Tag: "ipv6"})
	list = append(list, CustomValidate{Fn: CheckDigit, Tag: "digit"})
	list = append(list, CustomValidate{Fn: CheckPIN, Tag: "pin"})
	list = append(list, CustomValidate{Fn: CheckAbc, Tag: "abc"})
	list = append(list, CustomValidate{Fn: CheckId, Tag: "id"})
	list = append(list, CustomValidate{Fn: CheckCode, Tag: "code"})
	list = append(list, CustomValidate{Fn: CheckCountryCode, Tag: "country_code"})
	list = append(list, CustomValidate{Fn: CheckUsername, Tag: "username"})
	list = append(list, CustomValidate{Fn: CheckPattern, Tag: "regex"})
	list = append(list, CustomValidate{Fn: CheckAfterNow, Tag: "after_now"})
	list = append(list, CustomValidate{Fn: CheckNowOrAfter, Tag: "now_or_after"})
	return
}
func CheckString(fl validator.FieldLevel, fn func(string) bool) bool {
	s := fl.Field().String()
	if len(s) == 0 {
		return true
	}
	return fn(s)
}
func CheckEmail(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsEmail)
}
func CheckUrl(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsUrl)
}
func CheckUri(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsUri)
}
func CheckFax(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsFax)
}
func CheckPhone(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsPhone)
}
func CheckIp(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsIpAddress)
}
func CheckIpV4(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsIpAddressV4)
}
func CheckIpV6(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsIpAddressV6)
}
func CheckDigit(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsDigit)
}
func CheckPIN(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsDigit)
}
func CheckAbc(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsAbc)
}
func CheckId(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsCode)
}
func CheckCode(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsDashCode)
}
func CheckCountryCode(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsCountryCode)
}
func CheckUsername(fl validator.FieldLevel) bool {
	return CheckString(fl, s.IsUserName)
}
func CheckPattern(fl validator.FieldLevel) bool {
	param := fl.Param()
	if pattern, ok := PatternMap[param]; ok {
		return s.IsValidPattern(pattern, fl.Field().String())
	} else {
		panic("invalid pattern")
	}
}
// CheckAfterNow validates if the given time is greater than the current time
func CheckAfterNow(fl validator.FieldLevel) bool {
	var inputTime time.Time

	switch t := fl.Field().Interface().(type) {
	case string:
		parsedTime, err := time.Parse(time.RFC3339, t)
		if err != nil {
			return false
		}
		inputTime = parsedTime
	case time.Time:
		inputTime = t
	case *time.Time:
		inputTime = *t
	default:
		return false
	}

	return inputTime.UTC().After(time.Now().UTC())
}

// CheckNowOrAfter validates if the given time is greater or equal than the current time
func CheckNowOrAfter(fl validator.FieldLevel) bool {
	var inputTime time.Time

	switch t := fl.Field().Interface().(type) {
	case string:
		parsedTime, err := time.Parse(time.RFC3339, t)
		if err != nil {
			return false
		}
		inputTime = parsedTime
	case time.Time:
		inputTime = t
	case *time.Time:
		inputTime = *t
	default:
		return false
	}

	return inputTime.UTC().After(time.Now().UTC()) || inputTime.UTC().Equal(time.Now().UTC())
}
