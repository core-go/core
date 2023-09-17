package strings

import "strings"

func FixHeader(s string) string {
	s2 := strings.ToUpper(s)
	s1 := string(s2[0])
	for i := 1; i < len(s); i++ {
		if string(s[i-1]) == "-" || string(s[i-1]) == "_" {
			s1 = s1[:len(s1)-1]
			s1 += "-"
			s1 += string(s2[i])
		} else {
			s1 += string(s[i])
		}
	}
	return s1
}
func Fix(m map[string]string) map[string]string {
	rt := make(map[string]string)
	for k, v := range m {
		k2 := FixHeader(k)
		rt[k2] = v
	}
	return rt
}
func FixID(m map[string]string) map[string]string {
	rt := make(map[string]string)
	for k, v := range m {
		k2 := strings.Replace(k, "-id", "-ID", -1)
		k2 = FixHeader(k2)
		rt[k2] = v
	}
	return rt
}
