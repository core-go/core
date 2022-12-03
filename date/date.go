package date

import (
	"database/sql/driver"
	"strings"
	"time"
)

func Today() Date {
	today := time.Now().Truncate(24 * time.Hour)
	return Date(today)
}

type Date time.Time

func (d *Date) Time() time.Time {
	return time.Time(*d)
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.String() + "\""), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", value) //parse time
	if err != nil {
		return err
	}

	*d = Date(t)
	return nil
}

func (d *Date) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		*d = Date(vt.Truncate(24 * time.Hour))
	case string:
		tTime, err := time.Parse(time.RFC3339, vt)
		if err != nil {
			return err
		}
		*d = Date(tTime.Truncate(24 * time.Hour))
	case []byte:
		tTime, err := time.Parse(time.RFC3339, string(vt))
		if err != nil {
			return err
		}
		*d = Date(tTime.Truncate(24 * time.Hour))
	}
	return nil
}

func (d Date) Value() (driver.Value, error) {
	roundValue := d.Time().Truncate(time.Hour * 24)
	return roundValue, nil
}

func (d *Date) String() string {
	return d.Time().Format("2006-01-02")
}
