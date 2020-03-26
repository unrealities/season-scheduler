package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

type config struct {
	NumberOfTeams int `json:"numberOfTeams"`
}

func main() {
	flag.Parse()
	configFile := flag.String("config", "config.json", "location of config file.")
	config, err := parseConfig(configFile)
	if err != nil {
		fmt.Println("unable to parse config file")
	}
	fmt.Println(config.NumberOfTeams)
}

func parseConfig(file *string) (config, error) {
	var config config
	cfg, err := ioutil.ReadFile(*file)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(cfg, &config)

	return config, err
}
