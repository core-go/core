package template

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	TypeText       = "text"
	TypeIsNotEmpty = "isNotEmpty"
	TypeIsEmpty    = "isEmpty"
	TypeIsEqual    = "isEqual"
	TypeIsNotEqual = "isNotEqual"
	TypeIsNull     = "isNull"
	TypeIsNotNull  = "isNotNull"
	ParamText      = "text"
)

var ns = []string{"isNotNull", "isNull", "isEqual", "isNotEqual", "isEmpty", "isNotEmpty"}

func isValidNode(n string) bool {
	for _, s := range ns {
		if n == s {
			return true
		}
	}
	return false
}

type StringFormat struct {
	Texts      []string    `yaml:"" mapstructure:"texts" json:"texts,omitempty" gorm:"column:texts" bson:"texts,omitempty" dynamodbav:"texts,omitempty" firestore:"texts,omitempty"`
	Parameters []Parameter `yaml:"" mapstructure:"parameters" json:"parameters,omitempty" gorm:"column:parameters" bson:"parameters,omitempty" dynamodbav:"parameters,omitempty" firestore:"parameters,omitempty"`
}
type Parameter struct {
	Name string `yaml:"" mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Type string `yaml:"" mapstructure:"type" json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
}
type TemplateNode struct {
	Type      string       `yaml:"type" mapstructure:"type" json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
	Text      string       `yaml:"text" mapstructure:"text" json:"text,omitempty" gorm:"column:text" bson:"text,omitempty" dynamodbav:"text,omitempty" firestore:"text,omitempty"`
	Property  string       `yaml:"property" mapstructure:"property" json:"property,omitempty" gorm:"column:property" bson:"property,omitempty" dynamodbav:"property,omitempty" firestore:"property,omitempty"`
	Value     string       `yaml:"value" mapstructure:"value" json:"value,omitempty" gorm:"column:value" bson:"value,omitempty" dynamodbav:"value,omitempty" firestore:"value,omitempty"`
	Array     string       `yaml:"array" mapstructure:"array" json:"array,omitempty" gorm:"column:array" bson:"array,omitempty" dynamodbav:"array,omitempty" firestore:"array,omitempty"`
	Separator string       `yaml:"separator" mapstructure:"separator" json:"array,omitempty" gorm:"column:separator" bson:"separator,omitempty" dynamodbav:"separator,omitempty" firestore:"separator,omitempty"`
	Prefix    string       `yaml:"prefix" mapstructure:"prefix" json:"array,omitempty" gorm:"column:prefix" bson:"prefix,omitempty" dynamodbav:"prefix,omitempty" firestore:"prefix,omitempty"`
	Suffix    string       `yaml:"suffix" mapstructure:"suffix" json:"array,omitempty" gorm:"column:suffix" bson:"suffix,omitempty" dynamodbav:"suffix,omitempty" firestore:"suffix,omitempty"`
	Format    StringFormat `yaml:"format" mapstructure:"format" json:"format,omitempty" gorm:"column:format" bson:"format,omitempty" dynamodbav:"format,omitempty" firestore:"format,omitempty"`
}
type Template struct {
	Id        string         `yaml:"id" mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Text      string         `yaml:"text" mapstructure:"text" json:"text,omitempty" gorm:"column:text" bson:"text,omitempty" dynamodbav:"text,omitempty" firestore:"text,omitempty"`
	Templates []TemplateNode `yaml:"templates" mapstructure:"templates" json:"templates,omitempty" gorm:"column:templates" bson:"templates,omitempty" dynamodbav:"templates,omitempty" firestore:"templates,omitempty"`
}
type TStatement struct {
	Query  string        `yaml:"query" mapstructure:"query" json:"query,omitempty" gorm:"column:query" bson:"query,omitempty" dynamodbav:"query,omitempty" firestore:"query,omitempty"`
	Params []interface{} `yaml:"params" mapstructure:"params" json:"params,omitempty" gorm:"column:params" bson:"params,omitempty" dynamodbav:"params,omitempty" firestore:"params,omitempty"`
	Index  int           `yaml:"index" mapstructure:"index" json:"index,omitempty" gorm:"column:index" bson:"index,omitempty" dynamodbav:"index,omitempty" firestore:"index,omitempty"`
}

func LoadTemplates(trim func(string) string, files ...string) (map[string]*Template, error) {
	if len(files) == 0 {
		return loadTemplates(trim, "configs/query.xml")
	}
	return loadTemplates(trim, files...)
}
func loadTemplates(trim func(string) string, files ...string) (map[string]*Template, error) {
	l := len(files)
	f0, er0 := ReadFile(files[0])
	if er0 != nil {
		return nil, er0
	}
	if trim != nil {
		f0 = trim(f0)
	}
	templates, er0 := BuildTemplates(f0)
	if er0 != nil {
		return nil, er0
	}
	if l >= 2 {
		for i := 1; i < l; i++ {
			file, err := ReadFile(files[i])
			if err != nil {
				return templates, err
			}
			sub, er := BuildTemplates(file)
			if er0 != nil {
				return templates, er
			}
			for key, element := range sub {
				templates[key] = element
			}
		}
	}
	return templates, nil
}
func BuildTemplates(stream string) (map[string]*Template, error) {
	data := []byte(stream)
	buf := bytes.NewBuffer(data)
	dec := xml.NewDecoder(buf)
	ns := make([]TemplateNode, 0)
	ts := make(map[string]*Template)
	texts := make([]string, 0)
	start := false
	id := ""
	for {
		token, er0 := dec.Token()
		if token == nil {
			break
		}
		if er0 != nil {
			return nil, er0
		}
		switch element := token.(type) {
		case xml.CharData:
			if start == true {
				s := string([]byte(element))
				if !isEmptyNode(s) {
					n := TemplateNode{Type: "text", Text: s}
					texts = append(texts, s)
					n.Format = buildFormat(n.Text)
					ns = append(ns, n)
				}
			}
		case xml.EndElement:
			n := element.Name.Local
			if n == "select" || n == "insert" || n == "update" || n == "delete" {
				t := Template{Id: id}
				t.Text = strings.Join(texts, " ")
				t.Templates = ns
				ts[id] = &t
				ns = make([]TemplateNode, 0)
				start = false
			}
		case xml.StartElement:
			n := element.Name.Local
			if n == "select" || n == "insert" || n == "update" || n == "delete" {
				id = getValue(element.Attr, "id")
				texts = make([]string, 0)
				start = true
			} else {
				if element.Name.Local == "if" {
					test := getValue(element.Attr, "test")
					if len(test) > 0 {
						n := buildIf(test)
						if n != nil {
							n.Array = getValue(element.Attr, "array")
							n.Prefix = getValue(element.Attr, "prefix")
							n.Suffix = getValue(element.Attr, "suffix")
							n.Separator = getValue(element.Attr, "separator")
							sub, er1 := dec.Token()
							if er1 != nil {
								return nil, er1
							}
							switch inner := sub.(type) {
							case xml.CharData:
								s2 := string([]byte(inner))
								n.Text = s2
								n.Format = buildFormat(n.Text)
								texts = append(texts, s2)
							}
							ns = append(ns, *n)
						}
					}
				} else {
					if isEmptyNode(element.Name.Local) {
						property := getValue(element.Attr, "property")
						v := getValue(element.Attr, "value")
						array := getValue(element.Attr, "array")
						prefix := getValue(element.Attr, "prefix")
						suffix := getValue(element.Attr, "suffix")
						separator := getValue(element.Attr, "separator")
						n := TemplateNode{Type: element.Name.Local, Property: property, Value: v, Array: array, Prefix: prefix, Suffix: suffix, Separator: separator}
						sub, er1 := dec.Token()
						if er1 != nil {
							return nil, er1
						}
						switch inner := sub.(type) {
						case xml.CharData:
							s2 := string([]byte(inner))
							n.Text = s2
							n.Format = buildFormat(n.Text)
							texts = append(texts, s2)
						}
						ns = append(ns, n)
					}
				}
			}
		}
	}
	return ts, nil
}
func isEmptyNode(s string) bool {
	v := strings.Replace(s, "\n", " ", -1)
	v = strings.Replace(v, "\r", " ", -1)
	v = strings.TrimSpace(s)
	return len(v) == 0
}

func BuildTemplate(stream string) (*Template, error) {
	data := []byte(stream)
	buf := bytes.NewBuffer(data)
	dec := xml.NewDecoder(buf)
	ns := make([]TemplateNode, 0)
	texts := make([]string, 0)
	for {
		token, er0 := dec.Token()
		if token == nil {
			break
		}
		if er0 != nil {
			return nil, er0
		}
		switch element := token.(type) {
		case xml.CharData:
			s := string([]byte(element))
			if s != "\n" {
				n := TemplateNode{Type: "text", Text: s}
				texts = append(texts, s)
				n.Format = buildFormat(n.Text)
				ns = append(ns, n)
			}
		case xml.StartElement:
			if element.Name.Local == "if" {
				test := getValue(element.Attr, "test")
				if len(test) > 0 {
					n := buildIf(test)
					if n != nil {
						n.Array = getValue(element.Attr, "array")
						n.Prefix = getValue(element.Attr, "prefix")
						n.Suffix = getValue(element.Attr, "suffix")
						n.Separator = getValue(element.Attr, "separator")
						sub, er1 := dec.Token()
						if er1 != nil {
							return nil, er1
						}
						switch inner := sub.(type) {
						case xml.CharData:
							s2 := string([]byte(inner))
							n.Text = s2
							n.Format = buildFormat(n.Text)
							texts = append(texts, s2)
						}
						ns = append(ns, *n)
					}
				}
			} else {
				if isEmptyNode(element.Name.Local) {
					property := getValue(element.Attr, "property")
					v := getValue(element.Attr, "value")
					array := getValue(element.Attr, "array")
					prefix := getValue(element.Attr, "prefix")
					suffix := getValue(element.Attr, "suffix")
					separator := getValue(element.Attr, "separator")
					n := TemplateNode{Type: element.Name.Local, Property: property, Value: v, Array: array, Prefix: prefix, Suffix: suffix, Separator: separator}
					sub, er1 := dec.Token()
					if er1 != nil {
						return nil, er1
					}
					switch inner := sub.(type) {
					case xml.CharData:
						s2 := string([]byte(inner))
						n.Text = s2
						n.Format = buildFormat(n.Text)
						texts = append(texts, s2)
					}
					ns = append(ns, n)
				}
			}
		}
	}
	t := Template{}
	t.Text = strings.Join(texts, " ")
	t.Templates = ns
	return &t, nil
}
func getValue(attrs []xml.Attr, name string) string {
	if len(attrs) <= 0 {
		return ""
	}
	for _, attr := range attrs {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}
func buildFormat(str string) StringFormat {
	str2 := str
	str2b := str
	var str3 string
	texts := make([]string, 0)
	parameters := make([]Parameter, 0)
	var from, i, j int
	for {
		i = strings.Index(str2b, "{")
		if i >= 0 {
			str3 = str2b[i+1:]
			j = strings.Index(str3, "}")
			if j >= 0 {
				pro := str2b[i+1 : i+j+1]
				if isValidProperty(pro) {
					p := Parameter{}
					p.Name = pro
					if i >= 1 {
						var chr = string(str2b[i-1])
						if chr == "#" {
							texts = append(texts, str2[:from+i-1])
							p.Type = "param"
						} else if chr == "$" {
							texts = append(texts, str2[:from+i-1])
							p.Type = "text"
						} else {
							texts = append(texts, str2[:from+i])
							p.Type = "text"
						}
					} else {
						texts = append(texts, str2[:from+i])
						p.Type = "text"
					}
					parameters = append(parameters, p)
					from = from + i + j + 2
					str2 = str2[from:]
					str2b = str2
					from = 0
				} else {
					from = i + 1
					str2b = str2[i+1:]
				}
			} else {
				from = i + 1
				str2b = str2[from:]
			}
		} else {
			texts = append(texts, str2)
			break
		}
	}
	f := StringFormat{}
	f.Texts = texts
	f.Parameters = parameters
	return f
}
func RenderTemplateNodes(obj map[string]interface{}, templateNodes []TemplateNode) []TemplateNode {
	nodes := make([]TemplateNode, 0)
	for _, sub := range templateNodes {
		t := sub.Type
		if t == TypeText {
			nodes = append(nodes, sub)
		} else {
			attr := ValueOf(obj, sub.Property)
			if t == TypeIsNotNull {
				if attr != nil {
					vo := reflect.Indirect(reflect.ValueOf(attr))
					if vo.Kind() == reflect.Slice {
						if vo.Len() > 0 {
							nodes = append(nodes, sub)
						}
					} else {
						nodes = append(nodes, sub)
					}
				} else {
					vo := reflect.Indirect(reflect.ValueOf(attr))
					if vo.Kind() == reflect.Slice {
						if vo.Len() > 0 {
							nodes = append(nodes, sub)
						}
					}
				}
			} else if t == TypeIsNull {
				if attr == nil {
					nodes = append(nodes, sub)
				} else {
					vo := reflect.Indirect(reflect.ValueOf(attr))
					if vo.Kind() == reflect.Slice {
						if vo.Len() == 0 {
							nodes = append(nodes, sub)
						}
					}
				}
			} else if t == TypeIsEqual {
				if attr != nil {
					s := fmt.Sprintf("%v", attr)
					if sub.Value == s {
						nodes = append(nodes, sub)
					}
				}
			} else if t == TypeIsNotEqual {
				if attr != nil {
					s := fmt.Sprintf("%v", attr)
					if sub.Value != s {
						nodes = append(nodes, sub)
					}
				}
			} else if t == TypeIsEmpty {
				if attr != nil {
					s := fmt.Sprintf("%v", attr)
					if len(s) == 0 {
						nodes = append(nodes, sub)
					}
				}
			} else if t == TypeIsNotEmpty {
				if attr != nil {
					s := fmt.Sprintf("%v", attr)
					if len(s) > 0 {
						nodes = append(nodes, sub)
					}
				}
			}
		}
	}
	return nodes
}
func isValidProperty(v string) bool {
	var len = len(v) - 1
	for i := 0; i <= len; i++ {
		var chr = string(v[i])
		if !((chr >= "0" && chr <= "9") || (chr >= "A" && chr <= "Z") || (chr >= "a" && chr <= "z") || chr == "_" || chr == ".") {
			return false
		}
	}
	return true
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
func buildIf(t string) *TemplateNode {
	i := strings.Index(t, "!=")
	if i > 0 {
		s1 := strings.TrimSpace(t[0:i])
		s2 := strings.TrimSpace(t[i+2:])
		if len(s1) > 0 {
			if s2 == "null" {
				return &TemplateNode{Type: "isNotNull", Property: s1}
			} else {
				return &TemplateNode{Type: "isNotEqual", Property: s1, Value: trimQ(s2)}
			}
		}
	} else {
		i = strings.Index(t, "==")
		if i > 0 {
			s1 := strings.TrimSpace(t[0:i])
			s2 := strings.TrimSpace(t[i+2:])
			if len(s1) > 0 {
				if s2 == "null" {
					return &TemplateNode{Type: "isNull", Property: s1}
				} else {
					return &TemplateNode{Type: "isEqual", Property: s1, Value: trimQ(s2)}
				}
			}
		}
	}
	return nil
}
func trimQ(s string) string {
	if strings.HasPrefix(s, "'") {
		s = s[1:]
	} else if strings.HasPrefix(s, `"`) {
		s = s[1:]
	}
	if strings.HasSuffix(s, "'") {
		s = s[len(s)-1:]
	} else if strings.HasSuffix(s, `"`) {
		s = s[len(s)-1:]
	}
	return s
}

func ReadFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	text := string(content)
	return text, nil
}
func Q(s string) string {
	if !(strings.HasPrefix(s, "%") && strings.HasSuffix(s, "%")) {
		return "%" + s + "%"
	} else if strings.HasPrefix(s, "%") {
		return s + "%"
	} else if strings.HasSuffix(s, "%") {
		return "%" + s
	}
	return s
}
func Prefix(s string) string {
	if strings.HasSuffix(s, "%") {
		return s
	} else {
		return s + "%"
	}
}

const (
	t0 = "2006-01-02 15:04:05"
	t1 = "2006-01-02T15:04:05Z"
	t2 = "2006-01-02T15:04:05-0700"
	t3 = "2006-01-02T15:04:05.0000000-0700"

	l1 = len(t1)
	l2 = len(t2)
	l3 = len(t3)
)
func Merge(obj map[string]interface{}, format StringFormat, skipArray bool, separator string, prefix string, suffix string) string {
	results := make([]string, 0)
	parameters := format.Parameters
	if len(separator) > 0 && len(parameters) == 1 {
		p := ValueOf(obj, parameters[0].Name)
		vo := reflect.Indirect(reflect.ValueOf(p))
		if vo.Kind() == reflect.Slice {
			l := vo.Len()
			if l > 0 {
				strs := make([]string, 0)
				for i := 0; i < l; i++ {
					ts := Merge(obj, format, true, "", "", "")
					strs = append(strs, ts)
				}
				results = append(results, strings.Join(strs, separator))
				return prefix + strings.Join(results, "") + suffix
			}
		}
	}
	texts := format.Texts
	length := len(parameters)
	for i := 0; i < length; i++ {
		results = append(results, texts[i])
		p := ValueOf(obj, parameters[i].Name)
		if p != nil {
			if parameters[i].Type == ParamText {
				results = append(results, fmt.Sprintf("%v", p))
			} else {
				vo := reflect.Indirect(reflect.ValueOf(p))
				if vo.Kind() == reflect.Slice {
					l := vo.Len()
					if l > 0 {
						if skipArray {
							vx, _ :=GetDBValue(p, 2, "")
							results = append(results, vx)
						} else {
							sa := make([]string, 0)
							for i := 0; i < l; i++ {
								model := vo.Index(i).Addr()
								vx, _ :=GetDBValue(model.Interface(), 4, "")
								sa = append(sa, vx)
							}
							results = append(results, strings.Join(sa, ","))
						}
					}
				} else {
					vx, _ := GetDBValue(p, 2, "")
					results = append(results, vx)
				}
			}
		}
	}
	if len(texts[length]) > 0 {
		results = append(results, texts[length])
	}
	return prefix + strings.Join(results, "") + suffix
}
func Build(obj map[string]interface{}, template Template) string {
	results := make([]string, 0)
	renderNodes := RenderTemplateNodes(obj, template.Templates)
	for _, sub := range renderNodes {
		skipArray := sub.Array == "skip"
		s := Merge(obj, sub.Format, skipArray, sub.Separator, sub.Prefix, sub.Suffix)
		if len(s) > 0 {
			results = append(results, s)
		}
	}
	return strings.Join(results, "")
}
type QueryBuilder struct {
	Template  Template
	ModelType *reflect.Type
	Map       func(interface{}, *reflect.Type, ...func(string, reflect.Type) string) map[string]interface{}
	BuildSort func(string, reflect.Type) string
	Q         func(string) string
}
type Builder interface {
	BuildQuery(f interface{}) string
}

func UseQuery(isTemplate bool, query func(interface{}) string, id string, m map[string]*Template, modelType *reflect.Type, mp func(interface{}, *reflect.Type, ...func(string, reflect.Type) string) map[string]interface{}, buildSort func(string, reflect.Type) string, opts ...func(string) string) (func(interface{}) string, error) {
	if !isTemplate {
		return query, nil
	}
	b, err := NewQueryBuilder(id, m, modelType, mp, buildSort, opts...)
	if err != nil {
		return nil, err
	}
	return b.BuildQuery, nil
}
func UseQueryBuilder(isTemplate bool, builder Builder, id string, m map[string]*Template, modelType *reflect.Type, mp func(interface{}, *reflect.Type, ...func(string, reflect.Type) string) map[string]interface{}, buildSort func(string, reflect.Type) string, opts ...func(string) string) (Builder, error) {
	if !isTemplate {
		return builder, nil
	}
	return NewQueryBuilder(id, m, modelType, mp, buildSort, opts...)
}
func NewQueryBuilder(id string, m map[string]*Template, modelType *reflect.Type, mp func(interface{}, *reflect.Type, ...func(string, reflect.Type) string) map[string]interface{}, buildSort func(string, reflect.Type) string, opts ...func(string) string) (*QueryBuilder, error) {
	t, ok := m[id]
	if !ok || t == nil {
		return nil, errors.New("cannot get the template with id " + id)
	}
	var q func(string) string
	if len(opts) > 0 {
		q = opts[0]
	} else {
		q = Q
	}
	return &QueryBuilder{Template: *t, ModelType: modelType, Map: mp, BuildSort: buildSort, Q: q}, nil
}
func (b *QueryBuilder) BuildQuery(f interface{}) string {
	m := b.Map(f, b.ModelType, b.BuildSort)
	if b.Q != nil {
		q, ok := m["q"]
		if ok {
			s, ok := q.(string)
			if ok {
				m["q"] = b.Q(s)
			}
		}
	}
	return Build(m, b.Template)
}

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

func WrapString(v string) string {
	if strings.Index(v, `'`) >= 0 {
		return join(`'`, strings.Replace(v, "'", "''", -1), `'`)
	}
	return join(`'`, v, `'`)
}

func GetDBValue(v interface{}, scale int8, layoutTime string) (string, bool) {
	switch v.(type) {
	case string:
		s0 := v.(string)
		if len(s0) == 0 {
			return "''", true
		}

		return WrapString(s0), true
	case bool:
		b0 := v.(bool)
		if b0 {
			return "true", true
		} else {
			return "false", true
		}
	case int:
		return strconv.Itoa(v.(int)), true
	case int64:
		return strconv.FormatInt(v.(int64), 10), true
	case int32:
		return strconv.FormatInt(int64(v.(int32)), 10), true
	case big.Int:
		var z1 big.Int
		z1 = v.(big.Int)
		return z1.String(), true
	case float64:
		if scale >= 0 {
			mt := "%." + strconv.Itoa(int(scale)) + "f"
			return fmt.Sprintf(mt, v), true
		}
		return fmt.Sprintf("'%f'", v), true
	case time.Time:
		tf := v.(time.Time)
		if len(layoutTime) > 0 {
			f := tf.Format(layoutTime)
			return WrapString(f), true
		}
		f := tf.Format(t0)
		return WrapString(f), true
	case big.Float:
		n1 := v.(big.Float)
		if scale >= 0 {
			n2 := Round(n1, int(scale))
			return fmt.Sprintf("%v", &n2), true
		} else {
			return fmt.Sprintf("%v", &n1), true
		}
	case big.Rat:
		n1 := v.(big.Rat)
		if scale >= 0 {
			return RoundRat(n1, scale), true
		} else {
			return n1.String(), true
		}
	case float32:
		if scale >= 0 {
			mt := "%." + strconv.Itoa(int(scale)) + "f"
			return fmt.Sprintf(mt, v), true
		}
		return fmt.Sprintf("'%f'", v), true
	default:
		if scale >= 0 {
			v2 := reflect.ValueOf(v)
			if v2.Kind() == reflect.Ptr {
				v2 = v2.Elem()
			}
			if v2.NumField() == 1 {
				f := v2.Field(0)
				fv := f.Interface()
				k := f.Kind()
				if k == reflect.Ptr {
					if f.IsNil() {
						return "null", true
					} else {
						fv = reflect.Indirect(reflect.ValueOf(fv)).Interface()
						sv, ok := fv.(big.Float)
						if ok {
							return sv.Text('f', int(scale)), true
						} else {
							return "", false
						}
					}
				} else {
					sv, ok := fv.(big.Float)
					if ok {
						return sv.Text('f', int(scale)), true
					} else {
						return "", false
					}
				}
			} else {
				return "", false
			}
		} else {
			return "", false
		}
	}
	return "", false
}
func ParseDates(args []interface{}, dates []int) []interface{} {
	if args == nil || len(args) == 0 {
		return nil
	}
	if dates == nil || len(dates) == 0 {
		return args
	}
	res := append([]interface{}{}, args...)
	for _, d := range dates {
		if d >= len(args) {
			break
		}
		a := args[d]
		if s, ok := a.(string); ok {
			switch len(s) {
			case l1:
				t, err := time.Parse(t1, s)
				if err == nil {
					res[d] = t
				}
			case l2:
				t, err := time.Parse(t2, s)
				if err == nil {
					res[d] = t
				}
			case l3:
				t, err := time.Parse(t3, s)
				if err == nil {
					res[d] = t
				}
			}
		}
	}
	return res
}
func Round(num big.Float, scale int) big.Float {
	marshal, _ := num.MarshalText()
	var dot int
	for i, v := range marshal {
		if v == 46 {
			dot = i + 1
			break
		}
	}
	a := marshal[:dot]
	b := marshal[dot : dot+scale+1]
	c := b[:len(b)-1]

	if b[len(b)-1] >= 53 {
		c[len(c)-1] += 1
	}
	var r []byte
	r = append(r, a...)
	r = append(r, c...)
	num.UnmarshalText(r)
	return num
}
func RoundRat(rat big.Rat, scale int8) string {
	digits := int(math.Pow(float64(10), float64(scale)))
	floatNumString := rat.RatString()
	sl := strings.Split(floatNumString, "/")
	a := sl[0]
	b := sl[1]
	c, _ := strconv.Atoi(a)
	d, _ := strconv.Atoi(b)
	intNum := c / d
	surplus := c - d*intNum
	e := surplus * digits / d
	r := surplus * digits % d
	if r >= d/2 {
		e += 1
	}
	res := strconv.Itoa(intNum) + "." + strconv.Itoa(e)
	return res
}
