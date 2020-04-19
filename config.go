package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type config struct {
	EndDate       time.Time `json:"endDate"`
	DoubleHeaders bool      `json:"doubleHeaders"`
	NumGames      int       `json:"numGames"`
	StartDate     time.Time `json:"startDate"`
	Teams         []team    `json:"teams"`
}

// parseConfig reads the json config file and loads teams and other config settings
func parseConfig(file *string) (config, error) {
	var config config
	cfg, err := ioutil.ReadFile(*file)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(cfg, &config)

	return config, err
}
