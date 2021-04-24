package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/esutil"
)

type FieldLoader struct {
	client    *elasticsearch.Client
	indexName string
	name      string
}

func NewFieldLoader(db *elasticsearch.Client, indexName, name string) *FieldLoader {
	return &FieldLoader{
		client:    db,
		indexName: indexName,
		name:      name,
	}
}

func (l *FieldLoader) Values(ctx context.Context, ids []string) ([]string, error) {
	var array []string
	query := make(map[string]interface{})
	query["terms"] = map[string]interface{}{"_id": ids}
	query["boost"] = 1.0
	req := esapi.SearchRequest{
		Index:          []string{l.indexName},
		Body:           esutil.NewJSONReader(query),
		TrackTotalHits: true,
		Pretty:         true,
	}
	res, err := req.Do(ctx, l.client)
	if err != nil {
		return array, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return array, fmt.Errorf("response error")
	}

	var temp map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&temp)
	if err != nil {
		return array, err
	}

	hits := temp["hits"].(map[string]interface{})["hits"].([]interface{})
	result := make([]map[string]interface{}, 0)
	if err := json.NewDecoder(esutil.NewJSONReader(hits)).Decode(&result); err != nil {
		return array, err
	}
	for idx := range result {
		array = append(array, result[idx][l.name].(string))
	}
	return array, nil
}
