package utils

import (
	"fmt"
	"regexp"
)

func PregMatch(pattern string, content string) []string {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println(err.Error())
	}
	return reg.FindStringSubmatch(content)
}
