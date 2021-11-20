package mail

import "context"

type TemplateLoader interface {
	Load(ctx context.Context, id string) (string, string, error)
}
