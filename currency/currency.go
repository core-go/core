package currency

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math/big"
)

type Currency struct {
	big.Float
}
func (c Currency) MarshalJSON() (text []byte, err error)  {
	buff := []byte(c.String())
	return buff,nil

}

func (c *Currency) UnmarshalJSON(data []byte) error {
	err := c.UnmarshalText(data)
	if err != nil {
		return err
	}
	return nil
}

func (c Currency) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Currency) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}
