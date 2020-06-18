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
