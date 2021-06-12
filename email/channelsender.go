package email

type ChannelCodeSender struct {
	Channel chan string
}

func (ccs *ChannelCodeSender) SendCode(to string, code string) error {
	ccs.Channel <- code
	return nil
}
