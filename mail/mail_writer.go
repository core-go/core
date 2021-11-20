package mail

import (
	"context"
	"errors"
)

type MailWriter struct {
	Sender     SimpleMailSender
	Goroutines bool
}

func NewMailWriter(sender SimpleMailSender, goroutines bool) *MailWriter {
	return &MailWriter{Sender: sender, Goroutines: goroutines}
}

func (w *MailWriter) Write(ctx context.Context, model interface{}) error {
	mail, ok := model.(*SimpleMail)
	if !ok {
		return errors.New("input must be SimpleMail")
	} else {
		if w.Goroutines {
			go w.Sender.Send(*mail)
			return nil
		} else {
			return w.Sender.Send(*mail)
		}
	}
}
