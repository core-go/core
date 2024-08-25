package reader

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

type FileReader struct {
	FileName string
}

func NewFileReader(buildFileName func() string) (*FileReader, error) {
	fileName := buildFileName()
	if len(strings.TrimSpace(fileName)) == 0 {
		return nil, errors.New("file name cannot be empty")
	}
	return &FileReader{FileName: fileName}, nil
}

func (fr *FileReader) Read(next func(lines string, err error, numLine int) error) error {
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
		if len(line) != 0 {
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
