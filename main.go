package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

var config = flag.String("config", "config.json", "location of config file.")

type config struct {
	NumberOfTeams int `json: number_of_teams`
}

func main() {
	flag.Parse()
	cfg, err := ioutil.ReadFile(*config)
	if err != nil {
		fmt.Println("unable to locate config file")
	}

	var config config
	err = json.Unmarshal(cfg, &config)
	if err != nil {
		fmt.Println("config file not properly formatted")
	}
}
