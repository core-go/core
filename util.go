package service

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

func IsPointer(s interface{}) int {
	if reflect.ValueOf(s).Kind() == reflect.Ptr {
		return 1
	}
	return -1
}

func GetValue(model interface{}, fieldName string) (interface{}, error) {
	valueObject := reflect.Indirect(reflect.ValueOf(model))
	numField := valueObject.NumField()
	for i := 0; i < numField; i++ {
		if fieldName == valueObject.Type().Field(i).Name {
			return reflect.Indirect(valueObject).FieldByName(fieldName).Interface(), nil
		}
	}
	return nil, fmt.Errorf("Error no found field: " + fieldName)
}

func SetValue(model interface{}, index int, value interface{}) (interface{}, error) {
	valueModelObject := reflect.Indirect(reflect.ValueOf(model))
	if valueModelObject.Kind() == reflect.Ptr {
		valueModelObject = reflect.Indirect(valueModelObject)
	}

	valueModelObject.Field(index).Set(reflect.ValueOf(value))
	return model, nil
}

func FindNotIn(all []string, itemsNotIn []string) string {
	var result = ""
	for i := 1; i < len(itemsNotIn); i++ {
		if IndexOf(itemsNotIn[i], all) < 0 {
			return itemsNotIn[i]
		}
	}
	return result
}

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func RemoveUniCode(str string) string {
	str = strings.ToLower(str)
	str = regexp.MustCompile(`à|á|ạ|ả|ã|â|ầ|ấ|ậ|ẩ|ẫ|ă|ằ|ắ|ặ|ẳ|ẵ`).ReplaceAllString(str, "a")
	str = regexp.MustCompile(`è|é|ẹ|ẻ|ẽ|ê|ề|ế|ệ|ể|ễ`).ReplaceAllString(str, "e")
	str = regexp.MustCompile(`ì|í|ị|ỉ|ĩ`).ReplaceAllString(str, "i")
	str = regexp.MustCompile(`ò|ó|ọ|ỏ|õ|ô|ồ|ố|ộ|ổ|ỗ|ơ|ờ|ớ|ợ|ở|ỡ`).ReplaceAllString(str, "o")
	str = regexp.MustCompile(`ù|ú|ụ|ủ|ũ|ư|ừ|ứ|ự|ử|ữ`).ReplaceAllString(str, "u")
	str = regexp.MustCompile(`ỳ|ý|ỵ|ỷ|ỹ`).ReplaceAllString(str, "y")
	str = regexp.MustCompile(`đ`).ReplaceAllString(str, "d")
	str = regexp.MustCompile(`!|@|%|\^|\*|\(|\)|\+|\=|\<|\>|\?|\/|,|\.|\:|\;|\'|\"|\&|\#|\[|\]|~|\$|_`).ReplaceAllString(str, "-")
	//// Find and replace the special characters by -
	str = regexp.MustCompile(`-+-`).ReplaceAllString(str, "-")
	str = regexp.MustCompile(`^\-+|\-+$`).ReplaceAllString(str, "")
	//// trim - at the beginning and the and of this string
	return str
}

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

// Left left-pads the string with pad up to len runes
// len may be exceeded if
func PadLeft(str string, length int, pad string) string {
	return strings.Repeat(pad, length-len(str)) + str
}
func PadRight(str string, length int, pad string) string {
	return str + strings.Repeat(pad, length-len(str))
}
func Generate(length int) string {
	max := int(math.Pow(float64(10), float64(length))) - 1
	return PadLeft(strconv.Itoa(rand.Intn(max)), length, "0")
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
func Marshal(v interface{}) ([]byte, error) {
	b, ok1 := v.([]byte)
	if ok1 {
		return b, nil
	}
	s, ok2 := v.(string)
	if ok2 {
		return []byte(s), nil
	}
	return json.Marshal(v)
}
