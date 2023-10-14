package copier

import "encoding/json"

func Copy(src interface{}, des interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &des)
}
