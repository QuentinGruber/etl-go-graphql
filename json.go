package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JsonData []map[string]interface{}

func jsonToArrayOfObject(data []byte) DataTypeMapArray {
	jsonDataCleaned := strings.Replace(string(data), "\ufeff", "", -1)
	var jsonData JsonData
	json.Unmarshal([]byte(jsonDataCleaned), &jsonData)
	var result DataTypeMapArray
	for _, obj := range jsonData {
		entry := make(DataTypeMap, len(obj))
		for key, value := range obj {
			valueAsString := fmt.Sprintf("%v", value)
			dataType := DataType{
				rawData:  valueAsString,
				typeName: getTypeFromString(valueAsString),
			}
			entry[key] = dataType
		}
		result = append(result, entry)

	}
	return result
}
