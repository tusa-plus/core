package email

import (
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

var ErrInvalidPortFormat = fmt.Errorf("port format is invalid")
var ErrSMTPServerFailed = fmt.Errorf("SMTP server couldn't send code")

type CodeFormater interface {
	FormatCode(email string, code string) string
}

type EmailCodeSender struct {
	sender     string
	login      string
	password   string
	host       string
	port       string
	formatCode CodeFormater
}

func (ecs *EmailCodeSender) SendCode(to string, code string) error {
	m := gomail.NewMessage()

	body := ecs.formatCode.FormatCode(ecs.sender, code)
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
