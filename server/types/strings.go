package types

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

func NewRandomHexString(len int) (string, error) {
	bytes := make([]byte, len/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

var multipleSpacesPattern = regexp.MustCompile(" +")
var multipleBreaksPattern = regexp.MustCompile("\n+")

func CleanString(s string) string {
	if s == "" {
		return s
	}

	s = strings.TrimSpace(s)
	s = strings.Replace(s, "{{", "", -1)
	s = strings.Replace(s, "}}", "", -1)
	s = multipleSpacesPattern.ReplaceAllString(s, " ")
	return multipleBreaksPattern.ReplaceAllString(s, "\n")
}
