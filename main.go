package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	flag.Parse()
	configFile := flag.String("config", "config.json", "location of config file.")
	config, err := parseConfig(configFile)
	if err != nil {
		fmt.Printf("unable to parse config file: %s", err)
	}

	rand.Seed(time.Now().UnixNano())

	// Instantiate a league schedule with team's schedules
	numTeams := len(config.Teams)
	lgSchedule := make(schedules, numTeams, numTeams)
	lgGames := schedule{}
	lgGameID := 0
	teamAvailability := map[int]bool{}
	for i := 0; i < numTeams; i++ {
		teamAvailability[i] = true
	}

	// TODO:
	// We cannot just loop through the teams as it will always leave a team or two at the end with
	// not enough games to play. Also that is not how we'd want it to work.
	// We need to generate available teams. Can start with random but then we need to add
	// weighting to ensure divisional games are played and games are not played too close together
	totalExpectedGames := config.NumGames * (numTeams / 2)
	for lgGameID < totalExpectedGames {
		// Pick two teams for a series to play against each other
		lgGameID++
		homeTeam := findTeam(config.Teams, teamAvailability)
		teamAvailabilityWithoutHomeTeam := teamAvailability
		teamAvailabilityWithoutHomeTeam[homeTeam.ID] = false
		awayTeam := findTeam(config.Teams, teamAvailabilityWithoutHomeTeam)
		fmt.Printf("Home Team: %+v | Away team: %+v \n", homeTeam, awayTeam)
	}
	// Generate games
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

				// check to see if opponent still needs to play games
				// TODO: This will consume all games against a single opponent

				gamesRemaining := config.NumGames - len(lgSchedule[i])
				if gamesRemaining < config.SeriesMax {
					log.Printf("unable to schedule another series for %d", i)
					break
				}

				// try to prevent impossible series. If a team can play a series with exactly how many games
				// they have remaining, then we set that to the final series length.
				seriesLength := gamesRemaining
				if gamesRemaining > config.SeriesMax {
					seriesLength = randSeriesLength(config.SeriesMin, config.SeriesMax)
				}
				series := 0
				for (len(lgSchedule[j]) < config.NumGames) && (series < seriesLength) && (len(lgSchedule[i]) < config.NumGames) {
					lgGameID++
					series++
					// TODO: Handle dates. Don't allow two games in one day
					// TODO: Allow config for double-headers (still has to be same two teams)
					htNextGame := config.Teams[j].NextPlayableDate(config.StartDate, config.DoubleHeaders, lgSchedule[j], seriesLength)
					atNextGame := config.Teams[i].NextPlayableDate(config.StartDate, config.DoubleHeaders, lgSchedule[i], seriesLength)
					nextGame := maxTime(htNextGame, atNextGame)

					newGame := game{ID: lgGameID, AwayTeam: i, HomeTeam: j, Time: nextGame}
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
			fmt.Println(fmt.Sprintf("%d: %d (%d)", i, teamGames, config.NumGames))
		}
	}
}

func findTeam(teams []team, ta map[int]bool) team {
	teamID, err := randAvailableTeamID(ta)
	if err != nil {
		fmt.Printf("No available teams: %s \n", err)
		return team{}
	}
	return teams[teamID]
}

func maxTime(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}

func randSeriesLength(min, max int) int {
	rng := max - min
	if rng <= 0 {
		rng = 0
	}
	return rand.Intn(rng+1) + min
}

func randAvailableTeamID(ta map[int]bool) (int, error) {
	availableTeamIDs := []int{}
	for teamID, available := range ta {
		if available {
			availableTeamIDs = append(availableTeamIDs, teamID)
		}
	}
	numAvailableTeams := len(availableTeamIDs)
	if numAvailableTeams == 0 {
		return -1, fmt.Errorf("there are no available teams")
	}
	if numAvailableTeams == 1 {
		return availableTeamIDs[0], nil
	}

	return rand.Intn(numAvailableTeams), nil
}

func updateTeamAvailibility(team team, maxGames int, gp map[int]int, ta map[int]bool) map[int]bool {
	gameCount := team.IncrementGameCount()
	gp[team.ID] = gameCount
	if gameCount >= maxGames {
		ta[team.ID] = false
	}
	return ta
}
