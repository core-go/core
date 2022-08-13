package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type PasscodeRepository struct {
	collection *mongo.Collection
	passcodeName  string
	expiredAtName string
}

func NewPasscodeRepository(db *mongo.Database, collectionName string, options ...string) *PasscodeRepository {
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
	return &PasscodeRepository{db.Collection(collectionName), passcodeName, expiredAtName}
}

func (p *PasscodeRepository) Save(ctx context.Context, id string, passcode string, expiredAt time.Time) (int64, error) {
	pass := make(map[string]interface{})
	pass["_id"] = id
	pass[p.passcodeName] = passcode
	pass[p.expiredAtName] = expiredAt
	idQuery := bson.M{"_id": id}
	return UpsertOne(ctx, p.collection, idQuery, pass)
}

func (p *PasscodeRepository) Load(ctx context.Context, id string) (string, time.Time, error) {
	idQuery := bson.M{"_id": id}
	x := p.collection.FindOne(ctx, idQuery)
	er1 := x.Err()
	if er1 != nil {
		if strings.Compare(fmt.Sprint(er1), "mongo: no documents in result") == 0 {
			return "", time.Now().Add(-24 * time.Hour), nil
		}
		return "", time.Now().Add(-24 * time.Hour), er1
	}
	k, er3 := x.DecodeBytes()
	if er3 != nil {
		return "", time.Now().Add(-24 * time.Hour), er3
	}

	code := strings.Trim(k.Lookup(p.passcodeName).String(), "\"")
	expiredAt := k.Lookup(p.expiredAtName).Time()
	return code, expiredAt, nil
}

func (p *PasscodeRepository) Delete(ctx context.Context, id string) (int64, error) {
	idQuery := bson.M{"_id": id}
	return DeleteOne(ctx, p.collection, idQuery)
}
