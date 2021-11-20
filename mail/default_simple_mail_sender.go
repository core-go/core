package mail

type DefaultSimpleMailSender struct {
	MailSender MailSender
}

func NewSimpleMailSender(mailService MailSender) *DefaultSimpleMailSender {
	return &DefaultSimpleMailSender{MailSender: mailService}
}

func (s *DefaultSimpleMailSender) Send(m SimpleMail) error {
	var contents = make([]Content, len(m.Content))
	for i, content := range m.Content {
		contents[i] = content
	}
	mail := NewMailInit(m.From, m.Subject, m.To, m.Cc, contents...)
	return s.MailSender.Send(*mail)
}

func NewSimpleHtmlMail(mailFrom Email, subject string, mailTo []Email, cc *[]Email, htmlContent string) *SimpleMail {
	html := NewContent("text/html", htmlContent)
	s := SimpleMail{
		From:    mailFrom,
		To:      mailTo,
		Cc:      cc,
		Subject: subject,
		Content: []Content{*html},
	}
	return &s
}
