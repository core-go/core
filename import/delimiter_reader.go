package impt

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

type DelimiterFileReader struct {
	FileName  string
}

func NewDelimiterFileReader(buildFileName func() string) (*DelimiterFileReader,error) {
	fileName := buildFileName();
	if len(strings.TrimSpace(fileName)) == 0 {
		return nil, errors.New("file name cannot be empty")
	}
	return &DelimiterFileReader{FileName: fileName}, nil;
}

func (fr *DelimiterFileReader) Read(next func(lines string, err error, numLine int) error) error {
	file, err := os.Open(fr.FileName)
	if err != nil {
		err = errors.New("cannot open file")
		next("", err, 0)
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			err := next(line, nil, i)
			if err != nil {
				return err
			}
		}
		i++
	}
	next("", io.EOF, i)
	return nil
}
