package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// Config is struct for Configuration object from YAML
type Config struct {
	GeocodeURL         string `yaml:"geocodeURL"`
	GeocodePath        string `yaml:"geocodePath"`
	ConcurrentRoutines int    `yaml:"concurrentRoutines"`
}

// GenerateConfig will assign an unmarshalled JSON object to Configuration
func GenerateConfig() Config {
	file, err := filepath.Abs("./config.yml")

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return config
}
