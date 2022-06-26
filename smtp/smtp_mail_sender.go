package smtp

import (
	"errors"
	"github.com/core-go/core/mail"
)

type SmtpMailSender struct {
	Dialer *gomail.Dialer
}

func NewSmtpMailSender(config DialerConfig) *SmtpMailSender {
	m := new(SmtpMailSender)
	m.Dialer = gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	return m
}

func (m *SmtpMailSender) Send(mail mail.Mail) error {
	message := gomail.NewMessage()
	if mail.From.Address == "" {
		return errors.New("must have from field")
	} else if len(mail.From.Name) == 0 {
		message.SetHeader("From", mail.From.Address)
	} else {
		message.SetAddressHeader("From", mail.From.Address, mail.From.Name)
	}
	if mail.ReplyTo != nil {
		if len(mail.ReplyTo.Name) == 0 {
			message.SetHeader("Reply-To", mail.ReplyTo.Address)
		} else {
			message.SetAddressHeader("Reply-To", mail.ReplyTo.Address, mail.ReplyTo.Name)
		}
	}

	to := make(map[string][]string)
	toAddress := make([]string, 0)
	if mail.To[0].Address == "" {
		return errors.New("must have at least 1 receiver")
	}
	for _, i := range mail.To {
		toAddress = append(toAddress, message.FormatAddress(i.Address, i.Name))
	}
	to["To"] = toAddress
	message.SetHeaders(to)

	if mail.Cc != nil {
		cc := make(map[string][]string)
		ccAddress := make([]string, 0)
		for _, i2 := range *mail.Cc {
			ccAddress = append(ccAddress, message.FormatAddress(i2.Address, i2.Name))
		}
		cc["Cc"] = ccAddress
		message.SetHeaders(cc)
	}

	if mail.Bcc != nil {
		bcc := make(map[string][]string)
		bccAddress := make([]string, 0)
		for _, i3 := range *mail.Bcc {
			bccAddress = append(bccAddress, message.FormatAddress(i3.Address, i3.Name))
		}
		bcc["Bcc"] = bccAddress
		message.SetHeaders(bcc)
	}

	message.SetHeader("Subject", mail.Subject)
	if len(mail.Content) > 0 {
		var l = len(mail.Content) - 1
		for i := 0; i <= l; i++ {
			message.SetBody(mail.Content[i].Type, mail.Content[i].Value)
		}
	}
	return m.Dialer.DialAndSend(message)
}
