package impt

import (
	"bufio"
	"encoding/csv"
	"errors"
	"golang.org/x/text/encoding"
	"io"
	"os"
	"strings"
)

type FileReader struct {
	FileName  string
	Delimiter rune
	Decoder   *encoding.Decoder
}

func NewFileReader(buildFileName func() string) (*FileReader, error) {
	return NewDelimiterFileReader(buildFileName, ',')
}
func NewDelimiterFileReader(buildFileName func() string, delimiter rune, opts ...*encoding.Decoder) (*FileReader, error) {
	var decoder *encoding.Decoder
	if len(opts) > 0 && opts[0] != nil {
		decoder = opts[0]
	}
	var fr FileReader
	fileName := buildFileName()
	if len(strings.TrimSpace(fileName)) == 0 {
		return nil, errors.New("file name cannot be empty")
	}
	fr.FileName = fileName
	fr.Delimiter = delimiter
	fr.Decoder = decoder
	return &fr, nil
}

func (fr *FileReader) Read(next func(lines []string, err error, numLine int) error) error {
	file, err := os.Open(fr.FileName)
	if err != nil {
		err = errors.New("cannot open file")
		next(make([]string, 0), err, 0)
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan() {
		if scanner.Text() != "" {
			err := next([]string{scanner.Text()}, nil, i)
			if err != nil {
				return err
			}
		}
		i++
	}
	next([]string{}, io.EOF, i)
	return nil
}

func (fr *FileReader) ReadDelimiterFile(next func(lines []string, err error, numLine int) error) error {
	file, err := os.Open(fr.FileName)
	if err != nil {
		next(make([]string, 0), err, 0)
	}
	var r *csv.Reader
	if fr.Decoder != nil {
		r = csv.NewReader(fr.Decoder.Reader(file))
	} else {
		r = csv.NewReader(file)
	}

	if fr.Delimiter != 0 {
		r.Comma = fr.Delimiter
	}

	defer file.Close()
	i := 1
	for {
		record, err := r.Read()
		err2 := next(record, err, i)
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
