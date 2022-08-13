package sql

import "fmt"

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
func RemoveItem(slice []string, val string) []string {
	for i, item := range slice {
		if item == val {
			return RemoveIndex(slice, i)
		}
	}
	return slice
}
func QuoteByDriver(key, driver string) string {
	switch driver {
	case DriverMysql:
		return fmt.Sprintf("`%s`", key)
	case DriverMssql:
		return fmt.Sprintf(`[%s]`, key)
	default:
		return fmt.Sprintf(`"%s"`, key)
	}
}
func BuildResult(result int64, err error) (int64, error) {
	if err != nil {
		return result, err
	} else {
		return result, nil
	}
}
