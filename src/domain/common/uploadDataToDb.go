package common

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

func UploadCsvToDb(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New("No se ha podido leer el archivo.")
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
		// do something with read line
		fmt.Printf("%+v\n", rec)
	}
	return nil
}
