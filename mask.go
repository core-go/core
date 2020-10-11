package service

import (
	"sort"
	"strings"
)

func Include(vs []string, v string) bool {
	for _, s := range vs {
		if v == s {
			return true
		}
	}
	return false
}
func IncludeOfSort(vs []string, v string) bool {
	i := sort.SearchStrings(vs, v)
	if i >= 0 && vs[i] == v {
		return true
	}
	return false
}
func ValueOf(m interface{}, path string) interface{} {
	arr := strings.Split(path, ".")
	i := 0
	var c interface{}
	c = m
	l1 := len(arr) - 1
	for i < len(arr) {
		key := arr[i]
		m2, ok := c.(map[string]interface{})
		if ok {
			c = m2[key]
		}
		if !ok || i >= l1 {
			return c
		}
		i++
	}
	return c
}

func Mask(s string, start int, end int, mask string) string {
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start+end >= len(s) {
		return strings.Repeat(mask, len(s))
	}
	return s[:start] + strings.Repeat(mask, len(s)-start-end) + s[len(s)-end:]
}
func MaskMargin(s string, start int, end int, mask string) string {
	if start >= end {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start >= len(s) {
		return strings.Repeat(mask, len(s))
	}
	if end >= len(s) {
		return strings.Repeat(mask, start) + s[start:]
	}
	return strings.Repeat(mask, start) + s[start:end] + strings.Repeat(mask, len(s)-end)
}
