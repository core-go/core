package mapper

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"strings"
	"time"
)
const RFC3339Mili = "2006-01-02T15:04:05.000Z0700"

var cache = make(map[string]*template.Template, 0)
func GetTemplate(key string) (*template.Template, error) {
	t := cache[key]
	if t != nil {
		return t, nil
	}
	n, err := template.New("result").Parse(key)
	if err != nil {
		return nil, err
	}
	cache[key] = n
	return n, nil
}

func BuildParameters(s string, c *map[string]interface{}) (*map[string]interface{}, error) {
	return BuildParametersFromReader(strings.NewReader(s), c)
}
func BuildParametersFromReader(reader io.Reader, c *map[string]interface{}) (*map[string]interface{}, error) {
	now := time.Now()
	var p map[string]interface{}
	err := json.NewDecoder(reader).Decode(&p)
	if err != nil {
		return nil, err
	}
	if c != nil {
		for k, v := range *c {
			p[k] = v
		}
	}
	p["now_nano"] = now.Format(time.RFC3339Nano)
	p["now_iso"] = now.Format(time.RFC3339)
	p["now_mili"] = now.Format(RFC3339Mili)
	p["now"] = now
	p["now_unix"] = now.Unix()
	p["now_unix_nano"] = now.UnixNano()
	return &p, nil
}

func Build(template string, data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte(""))
	t, er1 := GetTemplate(template)
	if er1 != nil {
		return nil, er1
	}
	er2 := t.Execute(buf, data)
	if er2 != nil {
		return nil, er2
	}
	return buf.Bytes(), nil
}
