package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

var configFile = flag.String("config", "config.json", "location of config file.")

type config struct {
	NumberOfTeams int `json:"numberOfTeams"`
}

func main() {
	flag.Parse()
	cfg, err := ioutil.ReadFile(*configFile)
	if err != nil {
		fmt.Println("unable to locate config file")
	}

	var config config
	err = json.Unmarshal(cfg, &config)
	if err != nil {
		fmt.Println("config file not properly formatted")
	}
}
