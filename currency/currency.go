package currency

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math/big"
	"reflect"
)

type Currency struct {
	big.Float
}
func (c Currency) MarshalJSON() (text []byte, err error)  {
	buff := []byte(c.String())
	return buff, nil
}

func (c *Currency) UnmarshalJSON(data []byte) error {
	err := c.UnmarshalText(data)
	if err != nil {
		return err
	}
	return nil
}

func (c Currency) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *Currency) Scan(value interface{}) error {
	valType := reflect.TypeOf(value)
	if valType.Kind() == reflect.Slice {
		b, ok := value.([]byte)
		if !ok {
			return errors.New("type assertion to []byte failed")
		}
		return json.Unmarshal(b, &c)
	} else if valType.Kind() == reflect.String {
		b, ok := value.(string)
		if !ok {
			return errors.New("type assertion to string failed")
		}
		if b == "" {
			return nil
		}
		return json.Unmarshal([]byte(b), &c)
	}
	return nil
}
