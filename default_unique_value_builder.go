package service

import (
	"context"
	"fmt"
	"strings"
)

type DefaultUniqueValueBuilder struct {
	Generator Generator
	Loader    ValuesLoader
	Name      string
	Max       int
	ShortId   bool
}

func NewUniqueValueBuilder(generator Generator, loader ValuesLoader, name string, max int, shortId bool) *DefaultUniqueValueBuilder {
	return &DefaultUniqueValueBuilder{
		Generator: generator,
		Loader:    loader,
		Name:      name,
		Max:       max,
		ShortId:   shortId,
	}
}

// name is the field is used for create urlId
func (b *DefaultUniqueValueBuilder) Build(ctx context.Context, model interface{}, name string) (string, error) {
	var finalUrlId = ""

	var limitPreUrlId, err1 = GetValue(model, name)
	if err1 != nil {
		return "", err1
	}

	var limitPreUrlIdStr = ""
	if IsPointer(limitPreUrlId) == 1 {
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

	var urlIds, er3 = b.Loader.Values(ctx, array20ItemPattern)
	if er3 != nil {
		return "", er3
	}
	if len(urlIds) == 0 {
		finalUrlId = preUrlId
	} else {
		var urlIdNeed = FindNotIn(urlIds, array20ItemPattern)
		if urlIdNeed == "" {
			var randomId = ""
			if b.ShortId == false {
				randomId = RandomId()
			} else {
				randomId, er3 = ShortId()
				if er3 != nil {
					return "", er3
				}
			}
			finalUrlId = preUrlId + "-" + randomId
		} else {
			finalUrlId = urlIdNeed
		}
	}
	return finalUrlId, nil
}
