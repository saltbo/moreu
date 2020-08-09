package service

import (
	"fmt"

	"github.com/saltbo/moreu/pkg/mail"
)

var defaultNotify = NewNotify()

type Notify struct {
	mail *mail.Mail
}

func NewNotify() *Notify {
	return &Notify{mail: mail.New("", "", "")}
}

func (n *Notify) SendMail(subject, body, to string) error {
	return n.mail.Send(subject, body, to)
}

func SignupNotify(email, link string) error {
	template := `
       <h3>账户激活链接</h3>
       <p><a href="%s">点击此处重置密码</a></p>
		<p>如果您没有进行账号注册请忽略！</p>
       `
	body := fmt.Sprintf(template, link)
	return defaultNotify.SendMail("账号注册成功，请激活您的账户", body, email)
}

func PasswordResetNotify(email, link string) error {
	template := `
       <h3>密码重置链接</h3>
       <p><a href="%s">点击此处重置密码</a></p>
		<p>如果您没有申请重置密码请忽略！</p>
       `
	body := fmt.Sprintf(template, link)
	return defaultNotify.SendMail("密码重置申请", body, email)
}
