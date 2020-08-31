package client

import (
	"net/http"
)

const headerKeyUx = "X-Moreu-Ux"

type Context interface {
	GetHeader(key string) string
}

func InjectUx(req *http.Request, subject string) {
	req.Header.Set(headerKeyUx, subject)
}

func GetUx(c Context) string {
	return c.GetHeader(headerKeyUx)
}
