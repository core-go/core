package mail

type MailSender interface {
	Send(mail Mail) error
}
