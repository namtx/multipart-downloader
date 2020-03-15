package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
)

func Join(files []string, fileName string) error {
	var buf bytes.Buffer

	for _, file := range files {
		partFileName := filepath.Join("/tmp", file)

		b, err := ioutil.ReadFile(partFileName)
		if err != nil {
			return err
		}

		buf.Write(b)
	}

	err := ioutil.WriteFile(fileName, buf.Bytes(), 0644)

	return err
}
