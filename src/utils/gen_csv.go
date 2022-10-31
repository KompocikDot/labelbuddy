package utils

import (
	"strings"
)

func GenerateCSV(data [][]string) string {
	var result string
	for _, subdata := range data {
		result += strings.Join(subdata, ",") + "\n"
	}
	return result
}