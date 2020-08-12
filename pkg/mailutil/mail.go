package mailutil

import (
	"log"
	"net"
	"net/smtp"
)

var defaultMail *Mail

func Init(conf Config) {
	mail, err := NewMail(conf)
	if err != nil {
		log.Fatalln(err)
	}

	defaultMail = mail
}

func Send(subject, to, body string) error {
	return defaultMail.Send(subject, to, body)
}

type Config struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Mail struct {
	conf Config

	auth smtp.Auth
}

func NewMail(conf Config) (*Mail, error) {
	host, _, err := net.SplitHostPort(conf.Host)
	if err != nil {
		return nil, err
	}

	return &Mail{
		conf: conf,
		auth: smtp.PlainAuth("", conf.User, conf.Password, host),
	}, nil
}

func (m *Mail) Send(subject, to, body string) error {
	contentType := "Content-Type: text/html; charset=UTF-8"

	msg := []byte("To: " + to + "\r\nFrom: " + m.conf.User + "\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	return smtp.SendMail(m.conf.Host, m.auth, m.conf.User, []string{to}, msg)
}
