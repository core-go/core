package impt

import (
	"encoding/csv"
	"errors"
	"golang.org/x/text/encoding"
	"io"
	"os"
	"strings"
)

type FixedlengthFileReader struct {
	FileName  string
	Decoder   *encoding.Decoder
}

func NewFixedlengthFileReader(buildFileName func() string, opts ...*encoding.Decoder) (*FixedlengthFileReader,error) {
	fileName := buildFileName();
	if len(strings.TrimSpace(fileName)) == 0 {
		return nil, errors.New("file name cannot be empty")
	}
	var decoder *encoding.Decoder
	if len(opts) > 0 && opts[0] != nil {
		decoder = opts[0]
	}
	return &FixedlengthFileReader{
		FileName: fileName,
		Decoder: decoder,
	}, nil;
}

func(fr *FixedlengthFileReader) Read(next func(lines string, err error, numLine int) error) error {
	file, err := os.Open(fr.FileName)
	if err != nil {
		next("", err, 0)
	}
	var r *csv.Reader
	if fr.Decoder != nil {
		r = csv.NewReader(fr.Decoder.Reader(file))
	} else {
		r = csv.NewReader(file)
	}

	defer file.Close()
	i := 1
	for {
		record, err := r.Read()
		var err2 error
		if record == nil {
			err2 = next( "", err, i)
		} else {
			err2 = next( record[0], err, i)
		}
		if err2 != nil {
			return err2
		}
		if err == io.EOF {
			break
		}
		i++
	}
	return err
}
