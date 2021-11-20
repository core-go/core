package mail

import "context"

type DefaultTemplateLoader struct {
	Subject string
	Body    string
}

func NewTemplateLoader(subject string, body string) *DefaultTemplateLoader {
	return &DefaultTemplateLoader{subject, body}
}

func NewTemplateLoaderByConfig(c TemplateConfig) *DefaultTemplateLoader {
	return &DefaultTemplateLoader{c.Subject, c.Body}
}

func (s *DefaultTemplateLoader) Load(ctx context.Context, id string) (string, string, error) {
	return s.Subject, s.Body, nil
}
