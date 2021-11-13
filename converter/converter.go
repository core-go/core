package converter

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)
const layout = "2006-01-02"
func TimeToMilliseconds(time string) (int64, error) {
	var h, m, s int
	_, err := fmt.Sscanf(time, "%02d:%02d:%02d", &h, &m, &s)
	if err != nil {
		return 0, err
	}
	return int64(h * 3600000 + m * 60000 + s * 1000), nil
}
func DateToUnixTime(s string) (int64, error) {
	date, err := time.Parse(layout, s)
	if err != nil {
		return 0, err
	}
	return date.Unix() * 1000, nil
}
func DateToUnixNano(s string) (int64, error) {
	date, err := time.Parse(layout, s)
	if err != nil {
		return 0, err
	}
	return date.UnixNano(), nil
}
func UnixTime(m int64) string {
	dateUtc := time.Unix(0, m* 1000000)
	return dateUtc.Format("2006-01-02")
}
func MillisecondsToTimeString(milliseconds int) string {
	hourUint := 3600000 //60 * 60 * 1000 = 3600000
	minuteUint := 60000 //60 * 1000 = 60000
	secondUint := 1000
	hour := milliseconds / hourUint
	milliseconds = milliseconds % hourUint
	minute := milliseconds / minuteUint
	milliseconds = milliseconds % minuteUint
	second := milliseconds / secondUint
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}
func StringToAvroDate(date *string) (*int, error) {
	if date == nil {
		return nil, nil
	}
	d, err := time.Parse(layout, *date)
	if err != nil {
		return nil, err
	}
	i := int(d.Unix() / 86400)
	return &i, nil
}
func ToAvroDate(date *time.Time) *int {
	if date == nil {
		return nil
	}
	i := int(date.Unix() / 86400)
	return &i
}
func RoundFloat(num float64, slice int) float64 {
	c := math.Pow10(slice)
	result := math.Ceil(num*c) / c
	return result
}
func Round(num big.Float, scale int) big.Float {
	marshal, _ := num.MarshalText()
	var dot int
	for i, v := range marshal {
		if v == 46 {
			dot = i + 1
			break
		}
	}
	a := marshal[:dot]
	b := marshal[dot : dot+scale+1]
	c := b[:len(b)-1]

	if b[len(b)-1] >= 53 {
		c[len(c)-1] += 1
	}
	var r []byte
	r = append(r, a...)
	r = append(r, c...)
	num.UnmarshalText(r)
	return num
}
func RoundRat(rat big.Rat, scale int8) string {
	digits := int(math.Pow(float64(10), float64(scale)))
	floatNumString := rat.RatString()
	sl := strings.Split(floatNumString, "/")
	a := sl[0]
	b := sl[1]
	c, _ := strconv.Atoi(a)
	d, _ := strconv.Atoi(b)
	intNum := c / d
	surplus := c - d*intNum
	e := surplus * digits / d
	r := surplus * digits % d
	if r >= d/2 {
		e += 1
	}
	res := strconv.Itoa(intNum) + "." + strconv.Itoa(e)
	return res
}
