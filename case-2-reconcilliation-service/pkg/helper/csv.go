package helper

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func ParseCSV(csvPath string) (dataset []map[string]string, err error) {
	file, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return dataset, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	headers, err := reader.Read()
	if err != nil {
		fmt.Println("Error reading headers:", err)
		return dataset, err
	}

	dataset = []map[string]string{}

	// iterate csv records the the last record
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error reading records: %v", err)
		}

		newRecord := make(map[string]string)
		for i, value := range record {
			newRecord[headers[i]] = value
		}

		dataset = append(dataset, newRecord)
	}

	return dataset, nil
}
