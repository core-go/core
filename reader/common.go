package reader

import (
	"math/big"
	"reflect"
	"strconv"
	"time"
)

const DateLayout string = "2006-01-02 15:04:05 +0700 +07"

var funcMap = map[string]func(f reflect.Value, line string, format string, scale int) error{
	"string":     HandleString,
	"*string":    HandleString,
	"time.Time":  HandleTime,
	"*time.Time": HandleTime,
	"bool":       HandleBool,
	"*bool":      HandleBool,
	"int":        HandleInt,
	"*int":       HandleInt,
	"int64":      HandleInt,
	"*int64":     HandleInt64,
	"int32":      HandleInt32,
	"*int32":     HandleInt32,
	"big.Int":    HandleBigInt,
	"*big.Int":   HandleBigInt,
	"float64":    HandleFloat64,
	"*float64":   HandleFloat64,
	"big.Float":  HandleBigFloat,
	"*big.Float": HandleBigFloat,
}

func HandleUnknown(f reflect.Value, line string, format string, scale int) error {
	return nil
}
func HandleString(f reflect.Value, line string, format string, scale int) error {
	if f.Kind() == reflect.Ptr {
		f.Set(reflect.ValueOf(&line))
	} else {
		f.SetString(line)
	}
	return nil
}
func HandleTime(f reflect.Value, line string, format string, scale int) error {
	var fieldDate time.Time
	var err error
	if len(format) > 0 {
		fieldDate, err = time.Parse(format, line)
	} else {
		fieldDate, err = time.Parse(DateLayout, line)
	}
	if err != nil {
		return err
	}
	if f.Kind() == reflect.Ptr {
		f.Set(reflect.ValueOf(&fieldDate))
	} else {
		f.Set(reflect.ValueOf(fieldDate))
	}
	return nil
}
func HandleFloat64(f reflect.Value, line string, format string, scale int) error {
	floatValue, err := strconv.ParseFloat(line, 64)
	if err != nil {
		return err
	}
	if f.Kind() == reflect.Ptr {
		f.Set(reflect.ValueOf(&floatValue))
	} else {
		f.SetFloat(floatValue)
	}
	return nil
}
func HandleInt32(f reflect.Value, line string, format string, scale int) error {
	value, err := strconv.Atoi(line)
	if err != nil {
		return err
	}
	var i int32 = int32(value)
	if f.Kind() == reflect.Ptr {
		f.Set(reflect.ValueOf(&i))
	} else {
		f.Set(reflect.ValueOf(i))
	}
	return nil
}
func HandleInt64(f reflect.Value, line string, format string, scale int) error {
	value, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		return err
	}
	if f.Kind() == reflect.Ptr {
		f.Set(reflect.ValueOf(&value))
	} else {
		f.SetInt(value)
	}
	return nil
}
func HandleInt(f reflect.Value, line string, format string, scale int) error {
	value, err := strconv.Atoi(line)
	if err != nil {
		return err
	}
	if f.Kind() == reflect.Ptr {
		f.Set(reflect.ValueOf(&value))
	} else {
		f.Set(reflect.ValueOf(value))
	}
	return nil
}
func HandleBool(f reflect.Value, line string, format string, scale int) error {
	boolValue, err := strconv.ParseBool(line)
	if err != nil {
		return err
	}
	if f.Kind() == reflect.Ptr {
		f.Set(reflect.ValueOf(&boolValue))
	} else {
		f.SetBool(boolValue)
	}
	return nil
}
func HandleBigFloat(f reflect.Value, line string, format string, scale int) error {
	bf := new(big.Float)
	if bfv, ok := bf.SetString(line); ok {
		if scale >= 0 && bfv != nil {
			k := Round(*bf, scale)
			bf = &k
		}
		if f.Kind() == reflect.Ptr {
			f.Set(reflect.ValueOf(bfv))
		} else {
			if bfv != nil {
				f.Set(reflect.ValueOf(*bfv))
			}
		}
	}
	return nil
}
func HandleBigInt(f reflect.Value, line string, format string, scale int) error {
	bf := new(big.Int)
	if bfv, oki := bf.SetString(line, 10); oki {
		if f.Kind() == reflect.Ptr {
			f.Set(reflect.ValueOf(bfv))
		} else {
			if bfv != nil {
				f.Set(reflect.ValueOf(*bfv))
			}
		}
	}
	return nil
}
