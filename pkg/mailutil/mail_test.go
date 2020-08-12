package mailutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMail_Send(t *testing.T) {
	Init(Config{
		Host:     "smtpdm.aliyun.com:25",
		User:     "Moreu",
		Sender:   "no-reply@saltbo.fun",
		Password: "mG077URe2Wh9",
	})
	err := Send("title", "saltbo@foxmail.com", "body")
	assert.NoError(t, err)
}
