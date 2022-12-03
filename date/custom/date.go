package date

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

const format = "2006-01-02"

type Date struct {
	value time.Time
	valid bool
}

func (s *Date) Set(v time.Time) {
	s.value = v
	s.valid = true
}

func (s Date) Get() time.Time {
	return s.value
}

func (s Date) IsNull() bool {
	return !s.valid
}

func NewCustomDate() Date {
	return Date{}
}

func NewDate(v time.Time) Date {
	return Date{
		value: v,
		valid: true,
	}
}

func (s Date) MarshalJSON() ([]byte, error) {
	if s.IsNull() {
		return json.Marshal(nil)
	}

	return json.Marshal(s.Get().Format(format))
}

func (s *Date) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	parse, err := time.Parse(format, v)
	if err != nil {
		return err
	}
	*s = NewDate(parse)

	return nil
}

func (s *Date) Scan(value interface{}) error {
	v, ok := value.(time.Time)
	if ok {
		*s = NewDate(v)
	}
	return nil
}

func ConvertDate(val Date, layout string) string {
	if !val.IsNull() {
		return val.value.Format(layout)
	}
	return ""
}
