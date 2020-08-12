package service

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func ActivateLink(origin, email, token string) string {
	return fmt.Sprintf("%s/moreu/login/%s", origin, encodeToKey(email, token))
}

func PasswordRestLink(origin, email, token string) string {
	return fmt.Sprintf("%s/moreu/password_reset/%s", origin, encodeToKey(email, token))
}

var base64Encode = base64.URLEncoding.EncodeToString
var base64Decode = base64.URLEncoding.DecodeString

const moreuSplitKey = "|moreu|"

func encodeToKey(email, token string) string {
	return base64Encode([]byte(email + moreuSplitKey + token))
}

func decodeFromKey(key string) (email, token string) {
	bb, _ := base64Decode(key)
	sss := strings.Split(string(bb), moreuSplitKey)
	email = sss[0]
	token = sss[1]
	return
}
