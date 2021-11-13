package currency

import (
	"database/sql/driver"
	"encoding/json"
	"math/big"
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
