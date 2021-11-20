package mail

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type PasscodeSender struct {
	MailSender     SimpleMailSender
	From           Email
	TemplateLoader TemplateLoader
}

func NewPasscodeSender(mailSender SimpleMailSender, from Email, templateLoader TemplateLoader) *PasscodeSender {
	return &PasscodeSender{mailSender, from, templateLoader}
}

func truncatingSprintf(str string, args ...interface{}) string {
	n := strings.Count(str, "%s")
	if n > len(args) {
		n = len(args)
	}
	return fmt.Sprintf(str, args[0:n]...)
}

func (s *PasscodeSender) Send(ctx context.Context, to string, code string, expireAt time.Time, params interface{}) error {
	diff := expireAt.Sub(time.Now())
	strDiffMinutes := fmt.Sprintf("%.f", diff.Minutes())
	subject, template, err := s.TemplateLoader.Load(ctx, to)
	if err != nil {
		return err
	}

	content := truncatingSprintf(template,
		to, code, strDiffMinutes,
		to, code, strDiffMinutes,
		to, code, strDiffMinutes,
		to, code, strDiffMinutes,
		to, code, strDiffMinutes,
		to, code, strDiffMinutes,
		to, code, strDiffMinutes,
		to, code, strDiffMinutes,
		to, code, strDiffMinutes,
		to, code, strDiffMinutes)

	toMail := params.(string)
	mailTo := []Email{{Address: toMail}}
	mailData := NewSimpleHtmlMail(s.From, subject, mailTo, nil, content)
	return s.MailSender.Send(*mailData)
}
