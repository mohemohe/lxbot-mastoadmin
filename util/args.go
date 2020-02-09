package util

import (
	"strings"
)

func Equals(t string, args ...string) bool {
	a := strings.Fields(t)
	if len(a) != len(args) {
		return false
	}

	for i, _ := range args {
		if a[i] != args[i] {
			return false
		}
	}
	return true
}

func Prefix(t string, args ...string) bool {
	a := strings.Fields(t)
	if len(a) < len(args) {
		return false
	}

	for i, _ := range args {
		if a[i] != args[i] {
			return false
		}
	}
	return true
}
