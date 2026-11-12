package model

import "strings"

func IsPrivateKey(key string) bool {
	parts := strings.Split(key, "/")
	for i, part := range parts {
		if part == "private" && i < len(parts)-1 {
			return true
		}
	}
	return false
}
