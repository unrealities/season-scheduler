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
	ID       int       `json:"ID"`
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

	// Instantiate a league schedule with team's schedules
	numTeams := len(config.Teams)
	lgSchedule := make(schedules, numTeams, numTeams)
	lgGames := schedule{}
	lgGameID := 0

	// Generate games
	// TODO: track the day
	for i := range lgSchedule {
		// check to see if a given team's schedule has reached the required number of games
		if len(lgSchedule[i]) < config.NumGames {
			for j := range lgSchedule {
				// a team cannot play itself
				if i == j {
					continue
				}

				// make sure a team is not playing more games than they should
				if len(lgSchedule[i]) >= config.NumGames {
					break
				}

				// TODO: Handle dates. Don't allow two games in one day
				// TODO: Allow config for double-headers (still has to be same two teams)

				// check to see if opponent still needs to play games
				// TODO: This will consume all games against a single opponent
				for len(lgSchedule[j]) < config.NumGames {
					lgGameID++
					gameTime := config.StartDate
					newGame := game{ID: lgGameID, AwayTeam: i, HomeTeam: j, Time: gameTime}
					lgSchedule[j] = append(lgSchedule[j], newGame)
					lgSchedule[i] = append(lgSchedule[i], newGame)
					lgGames = append(lgGames, newGame)
				}
			}
		}
	}

	// Show scheduled games
	for _, l := range lgGames {
		fmt.Println(l.prettyPrint(config))
	}

	// Verify every team gets the required number of games
	for i := range lgSchedule {
		teamGames := len(lgSchedule[i])
		if teamGames != config.NumGames {
			fmt.Println(fmt.Sprintf("%d: %d", i, teamGames))
		}
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

func (g game) prettyPrint(c config) string {
	aTeamID := g.AwayTeam
	hTeamID := g.HomeTeam
	awayTeam := c.Teams[aTeamID]
	homeTeam := c.Teams[hTeamID]
	prettyAwayTeam := fmt.Sprintf("%s %s", awayTeam.Region, awayTeam.Name)
	prettyHomeTeam := fmt.Sprintf("%s %s", homeTeam.Region, homeTeam.Name)

	prettyDate := g.Time.Format("[Jan_2]")

	return fmt.Sprintf("%d | %s @ %s %s", g.ID, prettyAwayTeam, prettyHomeTeam, prettyDate)
}
