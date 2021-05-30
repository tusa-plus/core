package email

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/tusa-plus/core/utils"
	"gopkg.in/gomail.v2"
	"os"
	"strconv"
	"time"
)

const (
	codeLen = 6
)

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

func (ev *EmailVerification) SendCode(to string) error {
	code := ev.rndgen.NextString(codeLen)
	ev.storage.Set(to, []byte(code), ev.codeExpiration)
	err := ev.sender.SendCode(to, code)
	return err
}

func (ev *EmailVerification) VerifyCode(code string, email string) error {
	data, err := ev.storage.Get(code)
	if err != nil || string(data) != email {
		return fmt.Errorf("bad code")
	}
	return nil
}

type EmailCodeSender struct{}

func (ecs *EmailCodeSender) SendCode(to string, code string) error {
	m := gomail.NewMessage()

	m.SetBody("text/html", code)

	m.SetHeaders(map[string][]string{
		"From": {m.FormatAddress(ConfigSmtpSender, ConfigSmtpPort)},
		"To":   {to},
	})

	port, err := strconv.Atoi(ConfigSmtpPort)
	if err != nil {
		return err
	}
	d := gomail.NewPlainDialer(ConfigSmtpHost, port, ConfigSmtpLogin, ConfigSmtpPass)
	err = d.DialAndSend(m)
	return err
}

type ChannelCodeSender struct {
	channel chan string
}

func (ccs *ChannelCodeSender) SendCode(to string, code string) error {
	ccs.channel <- code
	return nil
}
