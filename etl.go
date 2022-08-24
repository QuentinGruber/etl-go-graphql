package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type DataType struct {
	rawData  string
	typeName string
}

type DataTypeMap map[string]DataType

type DataTypeMapArray []DataTypeMap

func convertToCamelCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		if i == 0 {
			words[i] = strings.ToLower(word)
		} else {
			words[i] = strings.Title(word)
		}
	}
	return strings.Join(words, "")
}

func getFileType(fileName string) string {
	return strings.Split(fileName, ".")[1]
}

type DataMap map[string]DataTypeMapArray

func getDataMap() DataMap {
	// make a list of all files inside ./data
	files, err := ioutil.ReadDir(config.DataFolderPath)
	if err != nil {
		fmt.Println(err)
	}
	filesNb := len(files)
	if filesNb == 0 {
		fmt.Println("No files found in ./data")
		os.Exit(0)
	}
	dataMap := make(DataMap, filesNb)
	var wg sync.WaitGroup
	wg.Add(filesNb)
	mutex := &sync.RWMutex{}
	for _, f := range files {
		// save each file in a map
		go func(fileName string) {
			fileType := getFileType(fileName)
			data, err := ioutil.ReadFile("./data/" + fileName)
			if err != nil {
				fmt.Println(err)
			}
			name := convertToCamelCase(strings.Split(fileName, ".")[0])
			switch fileType {
			case "csv":
				mutex.Lock()
				dataMap[name] = csvToArrayOfObject(data)
				mutex.Unlock()
			case "json":
				mutex.Lock()
				dataMap[name] = jsonToArrayOfObject(data)
				mutex.Unlock()
			default:
				fmt.Println("Unknown file type " + fileType)
				fmt.Println("File " + fileName + " will not be processed")
			}
			defer wg.Done()
		}(f.Name())
	}
	wg.Wait()
	return dataMap
}
