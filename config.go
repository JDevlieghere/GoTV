package main

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
	Queue     []Episode
}

func (config Configuration) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("DOWNLOAD DIR:"))
	buffer.WriteString(fmt.Sprintf("\n\t%s", config.Directory))
	buffer.WriteString(fmt.Sprintf("\nQUALITY:"))
	buffer.WriteString(fmt.Sprintf("\n\t%s", config.Quality))
	buffer.WriteString(fmt.Sprintf("\nSERIES:"))
	for _, s := range config.Series {
		buffer.WriteString(fmt.Sprintf("\n\t%s", s))
	}
	buffer.WriteString(fmt.Sprintf("\nQUEUE:"))
	for i, e := range config.Queue {
		buffer.WriteString(fmt.Sprintf("\n\t%d - %s", i, e))
	}
	return buffer.String()
}

func (config *Configuration) dequeue() *Episode {
	if len(config.Queue) == 0 {
		return nil
	}
	episode := config.Queue[0]
	config.Queue = config.Queue[1:]
	return &episode
}

func (config *Configuration) enqueue(episode *Episode) {
	if episode == nil {
		return
	}
	for _, e := range config.Queue {
		if e == *episode {
			return
		}
	}
	config.Queue = append(config.Queue, *episode)
}

func (config *Configuration) save() error {
	data, err := json.Marshal(*config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(CONF_FILE, data, 666)
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

func getConfig() Configuration {
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}
	return config
}
