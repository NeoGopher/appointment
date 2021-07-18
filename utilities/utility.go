package utilities

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func ParseToken(token string) (string, string, error) {
	var userName, userType string

	text, err := ParseCode(token)
	if err != nil {
		return userName, userType, err
	}

	details := strings.Split(text, "|")

	if len(details) != 2 {
		return userName, userType, fmt.Errorf("Error while parsing token")
	}

	userName = details[0]
	userType = details[1]

	return userName, userType, nil
}

func GetCode(data string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(data))
}

func ParseCode(code string) (string, error) {
	s, err := base64.RawStdEncoding.DecodeString(code)
	if err != nil {
		return string(s), err
	}

	return string(s), nil
}
