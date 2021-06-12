package email

import (
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

var ErrInvalidPortFormat = fmt.Errorf("port format is invalid")
var ErrSMTPServerFailed = fmt.Errorf("SMTP server couldn't send code")

func FormatCode(email string, code string) string {
	return code
}

type EmailCodeSender struct {
	sender     string
	login      string
	password   string
	host       string
	port       string
	formatCode func(string, string) string
}

func (ecs *EmailCodeSender) SendCode(to string, code string) error {
	m := gomail.NewMessage()

	body := ecs.formatCode(ecs.sender, code)
	m.SetBody("text/html", body)

	m.SetHeaders(map[string][]string{
		"From": {m.FormatAddress(ecs.sender, ecs.port)},
		"To":   {to},
	})

	portInt, err := strconv.Atoi(ecs.port)
	if err != nil {
		return ErrInvalidPortFormat
	}
	d := gomail.NewDialer(ecs.host, portInt, ecs.login, ecs.password)
	err = d.DialAndSend(m)
	if err != nil {
		return ErrSMTPServerFailed
	}
	return nil
}

func NewEmailCodeSender(sender, login, password, host, port string) EmailCodeSender {
	return EmailCodeSender{
		sender:     sender,
		login:      login,
		password:   password,
		host:       host,
		port:       port,
		formatCode: FormatCode,
	}
}
