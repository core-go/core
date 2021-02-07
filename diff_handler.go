package service

import (
	"context"
	"net/http"
	"reflect"
)

type DiffModelConfig struct {
	Id     string `mapstructure:"id" json:"id,omitempty" gorm:"column:id" bson:"_id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Origin string `mapstructure:"origin" json:"origin,omitempty" gorm:"column:origin" bson:"origin,omitempty" dynamodbav:"origin,omitempty" firestore:"origin,omitempty"`
	Value  string `mapstructure:"value" json:"value,omitempty" gorm:"column:value" bson:"value,omitempty" dynamodbav:"value,omitempty" firestore:"value,omitempty"`
	By     string `mapstructure:"by" json:"by,omitempty" gorm:"column:by" bson:"by,omitempty" dynamodbav:"by,omitempty" firestore:"by,omitempty"`
}
type DiffHandler struct {
	WriteLog    func(ctx context.Context, resource string, action string, success bool, desc string) error
	DiffService DiffService
	ModelType   reflect.Type
	IdNames     []string
	Indexes     map[string]int
	Offset      int
	Resource    string
	LogError    func(context.Context, string)
	Config      *DiffModelConfig
}
func NewDiffHandler(diffService DiffService, modelType reflect.Type, resource string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, config *DiffModelConfig, option ...int) *DiffHandler {
	offset := 1
	if len(option) == 1 {
		offset = option[0]
	}
	idNames := GetListFieldsTagJson(modelType)
	indexs := GetIndexes(modelType)
	return &DiffHandler{WriteLog: writeLog, DiffService: diffService, ModelType: modelType, IdNames: idNames, Indexes: indexs, Resource: resource, Offset: offset, Config: config, LogError: logError}
}
func NewDiffHandlerWithIds(diffService DiffService, modelType reflect.Type, resource string, logError func(context.Context, string), writeLog func(context.Context, string, string, bool, string) error, idNames []string, config *DiffModelConfig, option ...int) *DiffHandler {
	offset := 1
	if len(option) == 1 {
		offset = option[0]
	}
	if len(idNames) == 0 {
		idNames = GetListFieldsTagJson(modelType)
	}
	indexs := GetIndexes(modelType)
	return &DiffHandler{WriteLog: writeLog, DiffService: diffService, ModelType: modelType, IdNames: idNames, Indexes: indexs, Resource: resource, Offset: offset, Config: config, LogError: logError}
}

func (c *DiffHandler) Diff(w http.ResponseWriter, r *http.Request) {
	id, err := BuildId(r, c.ModelType, c.IdNames, c.Indexes, c.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		result, err := c.DiffService.Diff(r.Context(), id)
		if err != nil {
			Error(w, r, http.StatusInternalServerError, InternalServerError, c.LogError, c.Resource, "diff", err, c.WriteLog)
		} else {
			if c.Config == nil {
				Succeed(w, r, http.StatusOK, result, c.WriteLog, c.Resource, "diff")
			} else {
				m := make(map[string]interface{})
				if result.Id != nil {
					m[c.Config.Id] = result.Id
				}
				if result.Origin != nil {
					m[c.Config.Origin] = result.Origin
				}
				if result.Value != nil {
					m[c.Config.Value] = result.Value
				}
				if len(result.By) > 0 {
					m[c.Config.By] = result.By
				}
				Succeed(w, r, http.StatusOK, m, c.WriteLog, c.Resource, "diff")
			}
		}
	}
}
