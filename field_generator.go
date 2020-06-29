package service

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type FieldGenerator struct {
}

func (d *FieldGenerator) Generate(ctx context.Context, name string) (string, error) {
	name = strings.TrimSpace(name)
	name = regexp.MustCompile(`\s`).ReplaceAllString(name, "-")
	//preUrlId = RemoveAccents(preUrlId)
	name = RemoveUniCode(name)
	name = regexp.MustCompile(`[^\x00-\x7F]`).ReplaceAllString(name, "") // remove non-ascii character
	return name, nil
}

func (d *FieldGenerator) Array(ctx context.Context, name string) ([]string, error) {
	var array = make([]string, 20)
	array[0] = name
	for i := 1; i < 20; i++ {
		if i <= 9 {
			array[i] = name + "-" + strconv.Itoa(i)
		} else if i == 10 {
			array[i] = name + "-" + d.getDateYYMMDD()
		} else if i >= 11 && i < 20 {
			var a = i % 10
			array[i] = name + "-" + d.getDateYYMMDD() + "-" + strconv.Itoa(a)
		}
	}
	return array, nil
}

func (d *FieldGenerator) getDateYYMMDD() string { // format YYMMDD
	var newDateString = time.Now().Format("20060102")
	dateStr := []rune(newDateString)
	safeSubstring := string(dateStr[2:len(newDateString)])
	return safeSubstring
}

