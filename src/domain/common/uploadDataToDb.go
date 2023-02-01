package common

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

func ReadDataFromCsv(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	index := 0
	lines := make([][]string, 0)
	if err != nil {
		return lines, errors.New("No se ha podido leer el archivo.")
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// log.Fatal(err)
		}
		if index != 0 {
			lines = append(lines, rec)
		}
		index += 1
	}
	return lines, nil
}
