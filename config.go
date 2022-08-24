package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Configuration struct {
	RoleMap             RoleMap              `json:"roleMap"`
	DatasetRestricted   DatasetRestrictedMap `json:"datasetRestricted"`
	RestrictAllDatasets bool                 `json:"restrictAllDatasets"`
	Auth0Domain         string               `json:"auth0Domain"`
	ServerAddr          string               `json:"serverAddr"`
	CryptoKey           string               `json:"cryptoKey"`
	DataFolderPath      string               `json:"dataFolderPath"`
}

var config Configuration

func setupConfigurationGlobals() {
	file, _ := os.Open("config.json")
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue, &config)
}
