package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
)

type HealthChecker struct {
	name        string
	projectId   string
	opts        []option.ClientOption
	credentials []byte
}

func NewFirestoreHealthChecker(name string, projectId string, opts ...option.ClientOption) *HealthChecker {
	return &HealthChecker{projectId: projectId, name: name, opts: opts}
}

func NewHealthCheckerWithProjectId(projectId string, opts ...option.ClientOption) *HealthChecker {
	return NewFirestoreHealthChecker("firestore", projectId, opts...)
}
func NewHealthChecker(ctx context.Context, credentials []byte, projectId string, options ...string) (*HealthChecker, error) {
	var name string
	if len(options) > 0 && len(options[0]) > 0 {
		name = options[0]
	} else {
		name = "firestore"
	}
	opts := option.WithCredentialsJSON(credentials)
	creds, er2 := transport.Creds(ctx, opts)
	if er2 != nil {
		return nil, er2
	}
	if creds == nil {
		return nil, errors.New("error: credentials is nil")
	}
	return NewFirestoreHealthChecker(name, projectId, opts), nil
}
func (s HealthChecker) Name() string {
	return s.name
}

func (s HealthChecker) Check(ctx context.Context) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	client, err := firestore.NewClient(ctx, s.projectId, s.opts...)
	if err != nil {
		return res, err
	}
	defer client.Close()
	return res, nil
}

func (s *HealthChecker) Build(ctx context.Context, data map[string]interface{}, err error) map[string]interface{} {
	if err == nil {
		return data
	}
	if data == nil {
		data = make(map[string]interface{}, 0)
	}
	data["error"] = err.Error()
	return data
}
