package mail

import (
	"net"
	"net/smtp"
)

type Mail struct {
	host, user, password string
}

func New(host string, user string, password string) *Mail {
	return &Mail{host: host, user: user, password: password}
}

func (m *Mail) Send(subject, body, to string) error {
	host, _, err := net.SplitHostPort(m.host)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", m.user, m.password, host)
	contentType := "Content-Type: text/html; charset=UTF-8"

	msg := []byte("To: " + to + "\r\nFrom: " + m.user + "\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	return smtp.SendMail(m.host, auth, m.user, []string{to}, msg)
}
