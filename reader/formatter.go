package reader

import (
	"errors"
	"golang.org/x/text/encoding"
)

type FileType string

const (
	DelimiterType   FileType = "Delimiter"
	FixedlengthType FileType = "Fixedlength"
)

type Reader interface {
	Read(next func(lines string, err error, numLine int) error) error
}

func NewReader(buildFileName func() string, csvType FileType, opts ...*encoding.Decoder) (Reader, error) {
	if csvType == DelimiterType {
		return NewDelimiterFileReader(buildFileName)
	}
	if csvType == FixedlengthType {
		return NewFixedlengthFileReader(buildFileName, opts...)
	}
	return nil, errors.New("bad csv type")
}
