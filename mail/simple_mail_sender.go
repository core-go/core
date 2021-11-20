package mail

type SimpleMailSender interface {
	Send(mail SimpleMail) error
}
