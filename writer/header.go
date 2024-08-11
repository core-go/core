package writer

import (
	"bytes"
	"reflect"
	"strings"
)

func BuildHeader(modelType reflect.Type, delimiter string, tagName string) []byte {
	var b bytes.Buffer
	numFields := modelType.NumField()
	for i := 0; i < numFields; i++ {
		fieldName := modelType.Field(i).Tag.Get(tagName)
		if len(fieldName) == 0 {
			continue
		}
		b.WriteString(fieldName)
		if i != numFields-1 {
			b.WriteString(delimiter)
		}
	}
	b.WriteByte(byte('\n'))
	return b.Bytes()
}
func BuildHeaderText(modelType reflect.Type, delimiter string, tagName string) string {
	headers := make([]string, 0)
	numFields := modelType.NumField()
	for i := 0; i < numFields; i++ {
		fieldName := modelType.Field(i).Tag.Get(tagName)
		if len(fieldName) == 0 {
			continue
		}
		headers = append(headers, fieldName)
	}
	return strings.Join(headers, delimiter)
}
