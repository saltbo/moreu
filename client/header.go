package moreu

import "strconv"

const HeaderUserIdKey = "X-Moreu-Sub"

type Context interface {
	GetHeader(key string) string
}

func GetUserId(c Context) int64 {
	sub := c.GetHeader(HeaderUserIdKey)
	uid, _ := strconv.ParseInt(sub, 10, 64)
	return uid
}
