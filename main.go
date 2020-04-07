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
	NumGames  int       `json:"numGames"`
	StartDate time.Time `json:"startDate"`
	Teams     []team    `json:"teams"`
}

type team struct {
	Conference int    `json:"conference"`
	Divison    int    `json:"division"`
	ID         int    `json:"ID"`
	Name       string `json:"name"`
	Region     string `json:"region"`
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
	numTeams := len(config.Teams)
	lgSchedule := make(schedules, numTeams, numTeams)

	// TODO: This is the hard part. Need to actually make a schedule
	// Trying to do the "dumb" thing and start with the base case where each team
	// has to play one other team with no other qualifications
	for g := 0; g < config.NumGames; g++ {
		fmt.Printf("processing game: %d", g)
		for i, s := range lgSchedule {
			// check to see if a given team's schedule has reached the required number of games
			if len(s) < config.NumGames {
				for j, t := range lgSchedule {
					// a team cannot play itself
					if i == j {
						continue
					}

					// make sure a team is not playing more games than they should
					if len(s) >= config.NumGames {
						break
					}

					// check to see if opponent still needs to play games
					if len(t) < config.NumGames {
						newGame := game{
							AwayTeam: i,
							HomeTeam: j,
						}
						fmt.Printf("making a new game: %+v", newGame)
						s = append(s, newGame)
						t = append(t, newGame)
					}
				}
			}
		}
	}

	fmt.Printf("lgSchedule: %+v", lgSchedule)
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
