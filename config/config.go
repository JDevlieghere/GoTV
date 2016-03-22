package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	CONF_FILE = "config.json"
)

type Configuration struct {
	File      string
	Directory string
}

func (config Configuration) String() string {
	return fmt.Sprintf("File: %s\nDirectory: %s\n", config.File, config.Directory)
}

func readConfig() (Configuration, error) {
	file, err := ioutil.ReadFile(CONF_FILE)
	if err != nil {
		return Configuration{}, err
	}

	var config Configuration
	err = json.Unmarshal(file, &config)
	if err != nil {
		return Configuration{}, err
	}
	return config, nil
}

func storeConfig(config Configuration) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(CONF_FILE, data, 666)
}

func GetConfig() Configuration {
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func GetFile() string {
	config := GetConfig()
	return config.File
}

func GetDirectory() string {
	config := GetConfig()
	return config.Directory
}

func SetFile(file string) error {
	config := GetConfig()
	config.File = file
	return storeConfig(config)
}

func SetDirectory(directory string) error {
	config := GetConfig()
	config.Directory = directory
	return storeConfig(config)
}
