package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	CONF_FILE = "config.json"
)

type Configuration struct {
	Directory string
	Series    []string
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

func Get() Configuration {
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}
	return config
}
