package email

type ChannelCodeSender struct {
	channel chan string
}

func (ccs *ChannelCodeSender) SendCode(to string, code string) error {
	ccs.channel <- code
	return nil
}
