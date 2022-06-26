package sendgrid

import (
	"errors"
	m "github.com/core-go/core/mail"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailSender struct {
	Client *sendgrid.Client
}

func NewSendGridMailSender(apiKey string) *SendGridMailSender {
	m := new(SendGridMailSender)
	m.Client = sendgrid.NewSendClient(apiKey)
	return m
}

func (s *SendGridMailSender) Send(m m.Mail) error {
	mailSend := mail.NewV3Mail()
	from := mail.NewEmail(m.From.Name, m.From.Address)
	var tos []*mail.Email
	mailSend.SetFrom(from)
	mailSend.Subject = m.Subject
	l := len(m.Content)
	for i := 0; i < l; i++ {
		c := mail.NewContent(m.Content[i].Type, m.Content[i].Value)
		mailSend.AddContent(c)
	}
	p := mail.NewPersonalization()
	l = len(m.To)
	for i := 0; i < l; i++ {
		tos = append(tos, mail.NewEmail(m.To[i].Name, m.To[i].Address))
	}
	if tos == nil {
		return errors.New("must have at least 1 receiver")
	}
	ccs := m.Cc
	bccs := m.Bcc
	if ccs != nil {
		for _, i2 := range *ccs {
			p.AddCCs(mail.NewEmail(i2.Name, i2.Address))
		}
	}
	if bccs != nil {
		for _, i2 := range *bccs {
			p.AddBCCs(mail.NewEmail(i2.Name, i2.Address))
		}
	}
	p.AddTos(tos...)
	mailSend.AddPersonalizations(p)
	res, err := s.Client.Send(mailSend)
	if err != nil {
		return err
	}
	if res.StatusCode == 202 || res.StatusCode == 200 {
		return nil
	}
	return errors.New("error: " + res.Body)
}
