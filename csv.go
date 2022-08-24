package main

import (
	"bufio"
	"fmt"
	"strings"
)

func getSeparator(headerLine string) string {
	commaCount := strings.Count(headerLine, ",")
	semiColon := strings.Count(headerLine, ";")
	if commaCount > semiColon {
		return ","
	} else {
		return ";"
	}
}

func StringToLines(s string) (lines []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}

func csvToArrayOfObject(data []byte) DataTypeMapArray {
	// get the list of hearders from the first line
	lines, err := StringToLines(string(data))
	if err != nil {
		fmt.Println(err)
	}
	firstLine := strings.Replace(lines[0], "\ufeff", "", -1)
	separator := getSeparator(firstLine)
	headers := strings.Split(firstLine, separator)
	// for each line after the first create an entry in the array
	var result DataTypeMapArray
	for i := 1; i < len(lines); i++ {
		// create a map entry for each line
		entry := make(DataTypeMap, len(headers))
		// split the line into an array of values
		values := strings.Split(lines[i], separator)
		// for each header create an entry in the map
		for j := 0; j < len(headers); j++ {
			// create a data type object
			dataType := DataType{
				rawData:  values[j],
				typeName: getTypeFromString(values[j]),
			}
			entry[headers[j]] = dataType

		}
		// add the entry to the result array
		result = append(result, entry)
	}
	return result
}
