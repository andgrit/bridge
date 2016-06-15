package configuration

import (
	"io/ioutil"
	"encoding/json"
)

const CONFIGURATION_FILE = "./configuration.json"

type AppConfiguration struct {
	MongoIp string
	DatabaseName string
	DropDatabase bool
}

func DefaultConfiguration() (*AppConfiguration, error) {
	return ConfigurationFile(CONFIGURATION_FILE)
}

func ConfigurationFile(filename string) (*AppConfiguration, error)  {
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ConfigurationFromBytes(fileContents)
}

func ConfigurationFromBytes(jsonBytes []byte) (*AppConfiguration, error)  {
	var configuration AppConfiguration
	err := json.Unmarshal(jsonBytes, &configuration)
	return &configuration, err
}
