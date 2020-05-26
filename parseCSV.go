package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

// ParseCSV reads the given CSV and converts to Array of Address Structs for JSON Post Request
func ParseCSV() Records {
	file := os.Args[1]

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	rows := csv.NewReader(f)

	_, err1 := rows.Read()
	if err1 != nil {
		panic(err1)
	}

	lines, err2 := rows.ReadAll()
	if err2 != nil {
		panic(err2)
	}

	var attribute Attribute
	var attributes Attributes
	var records Records

	for _, each := range lines {

		objIDString, strErr := strconv.Atoi(each[0])
		if strErr != nil {
			panic(strErr)
		}
		attribute.ObjectID = objIDString
		attribute.Street = each[1]
		attribute.City = each[2]
		attribute.State = each[3]
		attribute.Zip = each[4]
		attributes.Attributes = attribute
		records.Records = append(records.Records, attributes)
	}

	return records
}

// SliceRecords Slices array of Attributes into an array of arrays
func SliceRecords(allRecords Records) []Records {
	var slicedRecords []Records

	chunkSize := 1000
	var recordLength int
	recordLength = len(allRecords.Records)

	for i := 0; i < recordLength; i += chunkSize {
		end := i + chunkSize

		if end > recordLength {
			end = recordLength
		}

		slicedRecords = append(slicedRecords, Records{allRecords.Records[i:end]})
	}

	return slicedRecords
}
