package currency

import (
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"math/big"
	"strconv"
)

type Currency struct {
	big.Float
}

func NewCurrency(val float64) Currency {
	return Currency{
		*big.NewFloat(val),
	}
}

func (c Currency) MarshalJSON() (text []byte, err error)  {
	return json.Marshal(c.Text('f', 2))
}

func (c *Currency) UnmarshalJSON(data []byte) error {
	err := c.UnmarshalText(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Currency) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	parsed, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return err
	}
	c.SetFloat64(parsed)
	return nil
}

func (c Currency) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *Currency) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if ok {
		return json.Unmarshal(b, &c)
	}
	b2, ok2 := value.(string)
	if ok2 {
		if b2 == "" {
			return nil
		}
		return json.Unmarshal([]byte(b2), &c)
	}
	return nil
}
