package bufio

import (
	"bufio"
	"os"
)

func ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []string{}, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	file.Close()
	return lines, nil
}

func LoadPemFile(filePath string) (string, error) {
	publicKeyFile, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	fileInfo, err := publicKeyFile.Stat()
	if err != nil {
		return "", err
	}
	size := fileInfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)
	if err != nil {
		return "", err
	}
	publicKeyFile.Close()
	return string(pembytes), nil
}
