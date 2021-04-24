package id

import (
	"context"
	"fmt"
	"strings"
)

type DefaultUniqueValueBuilder struct {
	Generator  Generator
	Values     func(ctx context.Context, ids []string) ([]string, error)
	Name       string
	Max        int
	GenerateId func(ctx context.Context) (string, error)
}

func NewUniqueValueBuilder(generator Generator, values func(context.Context, []string) ([]string, error), name string, max int, idGenerator func(ctx context.Context) (string, error)) *DefaultUniqueValueBuilder {
	return &DefaultUniqueValueBuilder{
		Generator:  generator,
		Values:     values,
		Name:       name,
		Max:        max,
		GenerateId: idGenerator,
	}
}

// Build name is the field is used for create urlId
func (b *DefaultUniqueValueBuilder) Build(ctx context.Context, model interface{}, name string) (string, error) {
	var finalUrlId = ""

	var limitPreUrlId, err1 = getValue(model, name)
	if err1 != nil {
		return "", err1
	}

	var limitPreUrlIdStr = ""
	if isPointer(limitPreUrlId) == 1 {
		if limitPreUrlId == nil {
			return "", fmt.Errorf("value of " + name + " cannot be nil")
		}
		limitPreUrlIdStr = *limitPreUrlId.(*string)
	} else {
		limitPreUrlIdStr = limitPreUrlId.(string)
	}
	limitPreUrlIdStr = strings.Trim(limitPreUrlIdStr, " ")
	if len(limitPreUrlIdStr) == 0 {
		return "", fmt.Errorf("value of " + name + " cannot be empty")
	}

	if len(limitPreUrlIdStr) > b.Max {
		limitPreUrlId = limitPreUrlIdStr[:b.Max]
	}
	var preUrlId, er1 = b.Generator.Generate(ctx, limitPreUrlIdStr)
	if er1 != nil {
		return "", er1
	}
	var array20ItemPattern, er2 = b.Generator.Array(ctx, preUrlId)
	if er2 != nil {
		return "", er2
	}

	var urlIds, er3 = b.Values(ctx, array20ItemPattern)
	if er3 != nil {
		return "", er3
	}
	if len(urlIds) == 0 {
		finalUrlId = preUrlId
	} else {
		var urlIdNeed = findNotIn(urlIds, array20ItemPattern)
		if urlIdNeed == "" {
			randomId, er4 := b.GenerateId(ctx)
			if er4 != nil {
				return "", er4
			}
			finalUrlId = preUrlId + "-" + randomId
		} else {
			finalUrlId = urlIdNeed
		}
	}
	return finalUrlId, nil
}
