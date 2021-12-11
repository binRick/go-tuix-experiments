package sexpect

import (
	"strings"
	"unicode"
)

func TabToSpace(input string) string {
	var result []string
	for _, i := range input {
		switch {
		case unicode.IsSpace(i):
			result = append(result, " ")
		case !unicode.IsSpace(i):
			result = append(result, string(i))
		}
	}
	return strings.Join(result, "")
}
