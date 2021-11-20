package mail

import (
	"context"
	"encoding/json"
)

type MQMailSender struct {
	Publish    func(ctx context.Context, data []byte, attributes map[string]string) (string, error)
	Goroutines bool
}

func NewMQMailSender(publish func(ctx context.Context, data []byte, attributes map[string]string) (string, error), goroutines bool) *MQMailSender {
	return &MQMailSender{Publish: publish, Goroutines: goroutines}
}
func (s *MQMailSender) Send(m SimpleMail) error {
	if s.Goroutines {
		go Publish(m, s.Publish)
		return nil
	} else {
		return Publish(m, s.Publish)
	}
}

func Publish(m SimpleMail, publish func(context.Context, []byte, map[string]string) (string, error)) error {
	b, er1 := json.Marshal(m)
	if er1 == nil {
		_, er2 := publish(context.TODO(), b, nil)
		return er2
	} else {
		return er1
	}
}
