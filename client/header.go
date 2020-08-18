package client

import (
	"net/http"
	"strconv"
)

const headerKeyUserId = "X-Moreu-Sub"

type Context interface {
	GetHeader(key string) string
}

func InjectUserId(req *http.Request, subject string) {
	req.Header.Set(headerKeyUserId, subject)
}

func GetUserId(c Context) int64 {
	sub := c.GetHeader(headerKeyUserId)
	uid, _ := strconv.ParseInt(sub, 10, 64)
	return uid
}
