package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"time"
)

type PasscodeRepository struct {
	collection    *firestore.CollectionRef
	passcodeName  string
	expiredAtName string
}
func NewPasscodeService(client *firestore.Client, collectionName string, options ...string) *PasscodeRepository {
	return NewPasscodeRepository(client, collectionName, options...)
}
func NewPasscodeRepository(client *firestore.Client, collectionName string, options ...string) *PasscodeRepository {
	var passcodeName, expiredAtName string
	if len(options) >= 1 && len(options[0]) > 0 {
		expiredAtName = options[0]
	} else {
		expiredAtName = "expiredAt"
	}
	if len(options) >= 2 && len(options[1]) > 0 {
		passcodeName = options[1]
	} else {
		passcodeName = "passcode"
	}
	return &PasscodeRepository{client.Collection(collectionName), passcodeName, expiredAtName}
}

func (s *PasscodeRepository) Save(ctx context.Context, id string, passcode string, expiredAt time.Time) (int64, error) {
	pass := make(map[string]interface{})
	pass[s.passcodeName] = passcode
	pass[s.expiredAtName] = expiredAt
	_, err := s.collection.Doc(id).Set(ctx, pass)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (s *PasscodeRepository) Load(ctx context.Context, id string) (string, time.Time, error) {
	doc, err := s.collection.Doc(id).Get(ctx)
	if err != nil {
		return "", time.Now(), err
	}
	code, err := doc.DataAt(s.passcodeName)
	if err != nil {
		return "", time.Now(), err
	}
	expiredAt, err := doc.DataAt(s.expiredAtName)
	if err != nil {
		return "", time.Now(), err
	}
	return code.(string), expiredAt.(time.Time), nil
}

func (s *PasscodeRepository) Delete(ctx context.Context, id string) (int64, error) {
	return DeleteOne(ctx, s.collection, id)
}
