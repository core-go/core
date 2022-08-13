package impt

import (
	"context"
	"errors"
	"reflect"
	"golang.org/x/text/encoding"
)

type csvType string
const (
	DelimiterType csvType = "Delimiter"
	FixedlengthType csvType = "Fixedlength"
)

type Formater interface {
	ToStruct(ctx context.Context, line string, res interface{}) (error)
}

func NewFormater(modelType reflect.Type, csvType csvType, opts... string) (Formater, error) {
	if csvType == DelimiterType {
		return NewDelimiterFormatter(modelType, opts...)
	}
	if csvType == FixedlengthType {
		return NewFixedLengthFormatter(modelType)
	}
	return nil, errors.New("Bad csv type")
}

type Reader interface {
	Read(next func(lines string, err error, numLine int) error) error
}

func NewReader(buildFileName func() string, csvType csvType, opts... *encoding.Decoder) (Reader, error) {
	if csvType == DelimiterType {
		return NewDelimiterFileReader(buildFileName)
	}
	if csvType == FixedlengthType {
		return NewFixedlengthFileReader(buildFileName, opts...)
	}
	return nil, errors.New("Bad csv type")
}
