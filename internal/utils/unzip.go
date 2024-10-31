package utils

import (
	"compress/gzip"
	"io"
)

func Unzip(body io.Reader) ([]byte, error) {
	gz, err := gzip.NewReader(body)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	ret, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
