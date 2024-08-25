package reader

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

type CSVReader struct {
	Comma    rune
	FileName string
}

func NewCSVReader(comma rune, buildFileName func() string) (*CSVReader, error) {
	fileName := buildFileName()
	if len(strings.TrimSpace(fileName)) == 0 {
		return nil, errors.New("file name cannot be empty")
	}
	return &CSVReader{Comma: comma, FileName: fileName}, nil
}

func (fr *CSVReader) Read(next func(record []string, err error, numLine int) error) error {
	file, err := os.Open(fr.FileName)
	if err != nil {
		err = errors.New("cannot open file")
		next(nil, err, 0)
		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = fr.Comma
	i := 1
	for {
		record, er1 := reader.Read()
		if er1 != nil {
			if er1.Error() == "EOF" {
				next(nil, io.EOF, i)
				return nil
			} else {
				return er1
			}
		}
		// Process each record
		er2 := next(record, nil, i)
		if er2 != nil {
			return er2
		}
		i++
	}
}
