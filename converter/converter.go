package converter

import (
	"fmt"
	"time"
)

func TimeToMilliseconds(sTime string) (int64, error) {
	var h, m, s int
	_, err := fmt.Sscanf(sTime, "%02d:%02d:%02d", &h, &m, &s)
	if err != nil {
		return 0, err
	}
	return int64(h * 3600000 + m * 60000 + s * 1000), nil
}
func DateToUnixTime(sDate string) (int64, error) {
	layout := "2006:01:02"
	date, err := time.Parse(layout, sDate)
	if err != nil {
		return 0, err
	}
	return date.Unix() * 1000, nil
}
func DateToUnixNano(sDate string) (int64, error) {
	layout := "2006:01:02"
	date, err := time.Parse(layout, sDate)
	if err != nil {
		return 0, err
	}
	return date.UnixNano(), nil
}
func UnixTime(milliseconds int64) string {
	dateUtc := time.Unix(0, milliseconds* 1000000)
	return dateUtc.Format("2006:01:02")
}
func MillisecondsToTimeString(timeMilisecond int) string {
	hourUint := 3600000 //60 * 60 * 1000 = 3600000
	minuteUint := 60000 //60 * 1000 = 60000
	secondUint := 1000
	hour := timeMilisecond / hourUint
	timeMilisecond = timeMilisecond % hourUint
	minute := timeMilisecond / minuteUint
	timeMilisecond = timeMilisecond % minuteUint
	second := timeMilisecond / secondUint
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}
