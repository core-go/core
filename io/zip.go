package io

import (
	"archive/zip"
	"bytes"
	"os"
)

func SaveToMemoryZip(fullName string, content string) (*bytes.Buffer, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)
	// Create a new zip archive.
	w := zip.NewWriter(buf)
	// Add a to the archive.
	f, err := w.Create(fullName)
	if err != nil {
		return buf, err
	}
	_, err = f.Write([]byte(content))
	if err != nil {
		return buf, err
	}
	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		return buf, err
	}
	return buf, nil
}
func SaveToZip(fullnamezipfile string, fileName, content string) error {
	newZipFile, err := os.Create(fullnamezipfile)
	if err != nil {
		return err
	}
	zipBuff, err := SaveToMemoryZip(fileName, content)
	if err != nil {
		return err
	}
	_, err = newZipFile.Write(zipBuff.Bytes())
	if err != nil {
		return err
	}
	err = newZipFile.Close()
	if err != nil {
		return err
	}
	return nil
}
func SaveFilesToMemoryZip(files []File) (*bytes.Buffer, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Add some files to the archive.
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			return buf, err
		}
		_, err = f.Write([]byte(file.Content))
		if err != nil {
			return buf, err
		}
	}

	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		return buf, err
	}
	return buf, nil
}
func SaveFilesToZip(fullnamezipfile string, files []File) error {
	newZipFile, err := os.Create(fullnamezipfile)
	if err != nil {
		return err
	}
	zipBuff, err := SaveFilesToMemoryZip(files)
	if err != nil {
		return err
	}
	_, err = newZipFile.Write(zipBuff.Bytes())
	if err != nil {
		return err
	}
	err = newZipFile.Close()
	if err != nil {
		return err
	}
	return nil
}
