package utils

import (
	"strings"
)

/*
DeleteEmpty deletes empty values from slice
*/
func DeleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if strings.TrimSpace(str) != "" && strings.TrimSpace(str) != " " && strings.TrimSpace(str) != "\n" && len(strings.TrimSpace(str)) > 0 {
			str = strings.Replace(str, "*.", "", -1)
			r = append(r, strings.TrimSpace(str))
		}
	}
	return r
}

func MergeMaps(mp1 map[string]interface{}, mp2 map[string]interface{}) map[string]interface{} {
	mp3 := make(map[string]interface{})

	for k, v := range mp1 {
		if _, ok := mp1[k]; ok {
			mp3[k] = v
		}
	}

	for k, v := range mp2 {
		if _, ok := mp2[k]; ok {
			mp3[k] = v
		}
	}
	return mp3
}

func Unique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}
