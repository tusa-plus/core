package email

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tusa-plus/core/utils"
	"gopkg.in/gomail.v2"
)

var ErrInvalidCode = fmt.Errorf("invalid code")
var ConfigSmtpSender = os.Getenv("EMAIL_SENDER")
var ConfigSmtpLogin = os.Getenv("EMAIL_USER")
var ConfigSmtpPass = os.Getenv("EMAIL_PASSWORD")
var ConfigSmtpHost = os.Getenv("EMAIL_HOST")
var ConfigSmtpPort = os.Getenv("EMAIL_PORT")

type CodeSender interface {
	SendCode(to string, code string) error
}

type EmailVerification struct {
	symbols        string
	storage        fiber.Storage
	codeExpiration time.Duration
	rndgen         *utils.RandomGenerator
	sender         CodeSender
}

func (ev *EmailVerification) SendCode(to string, symbols string, codeLen int) error {
	ev.rndgen = utils.NewRandomGenerator(symbols)
	code := ev.rndgen.NextString(codeLen)
	err := ev.storage.Set(to, []byte(code), ev.codeExpiration)
	if err != nil {
		return err
	}
	err = ev.sender.SendCode(to, code)
	return err
}

func (ev *EmailVerification) VerifyCode(code string, email string) error {
	data, err := ev.storage.Get(email)
	if err != nil {
		return err
	}
	if string(data) != code {
		return ErrInvalidCode
	}
	return nil
}

type EmailCodeSender struct{}

func (ecs *EmailCodeSender) EmailSend(sender, host, port, login, password, to, code string) error {
	m := gomail.NewMessage()

	m.SetBody("text/html", code)

	m.SetHeaders(map[string][]string{
		"From": {m.FormatAddress(sender, port)},
		"To":   {to},
	})

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return err
	}
	d := gomail.NewDialer(host, portInt, login, password)
	err = d.DialAndSend(m)
	return err
}

func (ecs *EmailCodeSender) SendCode(to string, code string) error {
	err := ecs.EmailSend(
		ConfigSmtpSender, ConfigSmtpHost, ConfigSmtpPort,
		ConfigSmtpLogin, ConfigSmtpPass, to, code,
	)
	return err
}

type ChannelCodeSender struct {
	channel chan string
}

func (ccs *ChannelCodeSender) SendCode(to string, code string) error {
	ccs.channel <- code
	return nil
}
