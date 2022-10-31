package utils

import (
	"strings"
	"time"
)

func RegexStrDatesToDates(dates [][]string) []string {
	var converted []string
	for _, match := range dates {
		strDate := strings.Replace(match[2], "\\/", "-", -1)
		time, _ := time.Parse("02-01-06", strDate)
		timeStr := time.Format("02.01.2006")
		converted = append(converted, timeStr)
	}
	return converted
}

func ReplaceSizeUnicodesToString(toReplace string) string {
	var replacements = []string{"½", "⅔", "⅓"}
	for index, unicode := range []string{"\\u00bd", "\\u2154", "\\u2153"} {
		toReplace = strings.Replace(toReplace, unicode, replacements[index], -1)
	}
	return toReplace
}