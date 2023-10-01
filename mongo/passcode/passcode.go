package passcode

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"reflect"
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

func DeleteOne(ctx context.Context, coll *mongo.Collection, query bson.M) (int64, error) {
	result, err := coll.DeleteOne(ctx, query)
	if result == nil {
		return 0, err
	}
	return result.DeletedCount, err
}
func Exist(ctx context.Context, collection *mongo.Collection, id interface{}, objectId bool) (bool, error) {
	query := bson.M{"_id": id}
	if objectId {
		objId, err := primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return false, err
		}
		query = bson.M{"_id": objId}
	}
	x := collection.FindOne(ctx, query)
	if x.Err() != nil {
		if fmt.Sprint(x.Err()) == "mongo: no documents in result" {
			return false, nil
		} else {
			return false, x.Err()
		}
	}
	return true, nil
}
func UpsertOne(ctx context.Context, collection *mongo.Collection, filter bson.M, model interface{}) (int64, error) {
	defaultObjID, _ := primitive.ObjectIDFromHex("000000000000")

	if idValue := filter["_id"]; idValue == "" || idValue == 0 || idValue == defaultObjID {
		return InsertOne(ctx, collection, model)
	} else {
		isExisted, err := Exist(ctx, collection, idValue, false)
		if err != nil {
			return 0, err
		}
		if isExisted {
			update := bson.M{
				"$set": model,
			}
			result := collection.FindOneAndUpdate(ctx, filter, update)
			if result.Err() != nil {
				if fmt.Sprint(result.Err()) == "mongo: no documents in result" {
					return 0, nil
				} else {
					return 0, result.Err()
				}
			}
			return 1, result.Err()
		} else {
			return InsertOne(ctx, collection, model)
		}
	}
}
func InsertOne(ctx context.Context, collection *mongo.Collection, model interface{}) (int64, error) {
	result, err := collection.InsertOne(ctx, model)
	if err != nil {
		errMsg := err.Error()
		if strings.Index(errMsg, "duplicate key error collection:") >= 0 {
			return 0, nil
		} else {
			return 0, err
		}
	} else {
		if idValue, ok := result.InsertedID.(primitive.ObjectID); ok {
			valueOfModel := reflect.Indirect(reflect.ValueOf(model))
			typeOfModel := valueOfModel.Type()
			idIndex, _, _ := FindIdField(typeOfModel)
			if idIndex != -1 {
				mapObjectIdToModel(idValue, valueOfModel, idIndex)
			}
		}
		return 1, err
	}
}
func FindIdField(modelType reflect.Type) (int, string, string) {
	return FindField(modelType, "_id")
}
func FindField(modelType reflect.Type, bsonName string) (int, string, string) {
	numField := modelType.NumField()
	for i := 0; i < numField; i++ {
		field := modelType.Field(i)
		bsonTag := field.Tag.Get("bson")
		tags := strings.Split(bsonTag, ",")
		json := field.Name
		if tag1, ok1 := field.Tag.Lookup("json"); ok1 {
			json = strings.Split(tag1, ",")[0]
		}
		for _, tag := range tags {
			if strings.TrimSpace(tag) == bsonName {
				return i, field.Name, json
			}
		}
	}
	return -1, "", ""
}
func mapObjectIdToModel(id primitive.ObjectID, valueOfModel reflect.Value, idIndex int) {
	switch reflect.Indirect(valueOfModel).Field(idIndex).Kind() {
	case reflect.String:
		if _, err := setValue(valueOfModel, idIndex, id.Hex()); err != nil {
			log.Println("Err: " + err.Error())
		}
		break
	default:
		if _, err := setValue(valueOfModel, idIndex, id); err != nil {
			log.Println("Err: " + err.Error())
		}
		break
	}
}
func setValue(model interface{}, index int, value interface{}) (interface{}, error) {
	vo := reflect.Indirect(reflect.ValueOf(model))
	switch reflect.ValueOf(model).Kind() {
	case reflect.Ptr:
		{
			vo.Field(index).Set(reflect.ValueOf(value))
			return model, nil
		}
	default:
		if modelWithTypeValue, ok := model.(reflect.Value); ok {
			_, err := setValueWithTypeValue(modelWithTypeValue, index, value)
			return modelWithTypeValue.Interface(), err
		}
	}
	return model, nil
}
func setValueWithTypeValue(model reflect.Value, index int, value interface{}) (reflect.Value, error) {
	trueValue := reflect.Indirect(model)
	switch trueValue.Kind() {
	case reflect.Struct:
		{
			val := reflect.Indirect(reflect.ValueOf(value))
			if trueValue.Field(index).Kind() == val.Kind() {
				trueValue.Field(index).Set(reflect.ValueOf(value))
				return trueValue, nil
			} else {
				return trueValue, fmt.Errorf("value's kind must same as field's kind")
			}
		}
	default:
		return trueValue, nil
	}
}
