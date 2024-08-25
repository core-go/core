package writer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	writers = make(map[string]io.WriteCloser)
	mux     = &sync.Mutex{}
)

type FileWriter struct {
	FileName         string
	GenerateFileName func() string
	Out              *bufio.Writer
	File             *os.File
}

func NewFileWriter(buildFileName func() string) (*FileWriter, error) {
	fileName := buildFileName()
	dirPath := filepath.Dir(fileName)
	_, err := MakeDir(dirPath)
	if err != nil {
		return nil, err
	}

	outFile := buildFileName()
	f, err := OpenFileWriter(outFile)
	if err != nil {
		return nil, err
	}

	var fw FileWriter
	mux.Lock()
	defer mux.Unlock()
	out := bufio.NewWriter(f)
	fw.FileName = fileName
	fw.Out = out
	fw.File = f

	return &fw, nil
}

func (fw *FileWriter) Write(p []byte) (n int, err error) {
	mux.Lock()
	defer mux.Unlock()
	return fw.Out.Write(p)
}

// CloseWriter close the underlying writer.
func (fw *FileWriter) Close() error {
	fw.Out.Flush()
	return fw.File.Close()
}

func CloseAllWriters() error {
	mux.Lock()
	defer mux.Unlock()
	for _, w := range writers {
		if w != nil {
			w.Close()
		}
	}
	return nil
}

func DeleteFile(outFile string) (bool, error) {
	_, err := os.Stat(outFile)
	if !os.IsNotExist(err) {
		if err := os.Remove(outFile); err != nil {
			return true, err
		}
	}
	return false, nil
}

func AppendWriter(path string) bool {
	mux.Lock()
	_, ok := writers[path]
	mux.Unlock()
	return ok
}

func OpenFileWriter(outFile string) (*os.File, error) {
	if AppendWriter(outFile) {
		return nil, fmt.Errorf("log writer already opened on %s", outFile)
	}

	return os.OpenFile(outFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
}

func MakeDir(dirName string) (bool, error) {
	_, err := os.Stat(dirName)
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return true, err
		}
		if !info.IsDir() {
			return true, errors.New("path exists but is not a directory")
		}
		return true, nil
	} else {
		if err := os.MkdirAll(dirName, 0755); err != nil {
			return false, err
		}
	}
	return false, nil
}
