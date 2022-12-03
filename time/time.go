package datetime

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

const format = "2006-01-02T15:04:05.000000-07:00"

type DateTime struct {
	value time.Time
	valid bool
}

func (s *DateTime) Set(v time.Time) {
	s.value = v
	s.valid = true
}

func (s DateTime) Get() time.Time {
	return s.value
}

func (s DateTime) IsNull() bool {
	return !s.valid
}

func NewTime() DateTime {
	return DateTime{}
}

func NewDateTime(v time.Time) DateTime {
	return DateTime{
		value: v,
		valid: true,
	}
}

func (s DateTime) MarshalJSON() ([]byte, error) {
	if s.IsNull() {
		return json.Marshal(nil)
	}
	return json.Marshal(s.Get().Format(format))
}

func (s *DateTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	parse, err := time.Parse(format, v)
	if err != nil {
		return err
	}
	*s = NewDateTime(parse)

	return nil
}

func ConvertDateTime(val DateTime, layout string) string {
	if !val.IsNull() {
		return val.value.Format(layout)
	}
	return ""
}
