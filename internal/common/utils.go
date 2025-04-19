package common

import (
	"regexp"
	"strconv"
	"strings"
)

func StringAsContentString(raw string) (bool, string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return false, ""
	}
	return true, raw
}

func StringAsNumberOrNil(raw string) (bool, int) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return false, 0
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return false, 0
	}
	return true, value
}

func StringAsRegExpOrNil(raw string) (bool, *regexp.Regexp) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return false, nil
	}
	value, err := regexp.Compile(raw)
	if err != nil {
		return false, nil
	}
	return true, value
}
