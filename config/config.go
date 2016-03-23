package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	CONF_FILE = "config.json"
)

type Configuration struct {
	Directory string
	Quality   string
	Series    []string
}

func (config Configuration) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("DOWNLOAD DIR:"))
	buffer.WriteString(fmt.Sprintf("\n\t%s\n", config.Directory))
	buffer.WriteString(fmt.Sprintf("QUALITY:"))
	buffer.WriteString(fmt.Sprintf("\n\t%s\n", config.Quality))
	buffer.WriteString(fmt.Sprintf("SERIES:"))
	for _, s := range config.Series {
		buffer.WriteString(fmt.Sprintf("\n\t%s", s))
	}
	return buffer.String()
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
