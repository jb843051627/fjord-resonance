package export

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
)

func ZipCSV(name string, content []byte) ([]byte, error) {
	if name == "" {
		return nil, fmt.Errorf("archive name is empty")
	}
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	file, err := writer.Create(name)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(file, bytes.NewReader(content)); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
