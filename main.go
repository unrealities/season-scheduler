package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"time"
)

type config struct {
	EndDate   time.Time `json:"endDate"`
	StartDate time.Time `json:"startDate"`
	Teams     []team    `json:"teams"`
}

type team struct {
	Conference int    `json:"conference"`
	Divison    int    `json:"division"`
	ID         int    `json:"ID"`
	Name       string `json:"name"`
	Region     string `json:"regions"`
}

type game struct {
	AwayTeam int       `json:"awayTeam"`
	HomeTeam int       `json:"homeTeam"`
	Time     time.Time `json:"time"`
}

type schedule []game
type schedules []schedule

func main() {
	flag.Parse()
	configFile := flag.String("config", "config.json", "location of config file.")
	config, err := parseConfig(configFile)
	if err != nil {
		fmt.Printf("unable to parse config file: %s", err)
	}
	fmt.Printf("config: %+v", config)
	fmt.Printf("teams: %+v", config.Teams)

	// Instantiate a league schedule with team's schedules
	lgSchedule := schedules{}
	for i := range config.Teams {
		lgSchedule[i] = schedule{}
	}

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
