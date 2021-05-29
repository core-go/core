package strings

import "strings"

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
