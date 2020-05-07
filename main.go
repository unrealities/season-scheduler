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
		homeTeam := findTeam(config.Teams, teamAvailability)
		homeTeamGamesRemaining := config.NumGames - len(lgSchedule[homeTeam.ID])
		if homeTeamGamesRemaining < config.SeriesMax {
			log.Printf("unable to schedule another series for %+v", homeTeam)
			continue
		}

		seriesLength := homeTeamGamesRemaining
		if homeTeamGamesRemaining > config.SeriesMax {
			seriesLength = randSeriesLength(config.SeriesMin, config.SeriesMax)
		}
		if homeTeamGamesRemaining > config.SeriesMax {
			seriesLength = randSeriesLength(config.SeriesMin, config.SeriesMax)
		}
		series := 0

		teamAvailabilityWithoutHomeTeam := teamAvailability
		teamAvailabilityWithoutHomeTeam[homeTeam.ID] = false
		awayTeam := findTeam(config.Teams, teamAvailabilityWithoutHomeTeam)

		for (len(lgSchedule[homeTeam.ID]) < config.NumGames) && (series < seriesLength) && (len(lgSchedule[awayTeam.ID]) < config.NumGames) {
			lgGameID++
			series++
			// TODO: Handle dates. Don't allow two games in one day
			// TODO: Allow config for double-headers (still has to be same two teams)
			htNextGame := homeTeam.NextPlayableDate(config.StartDate, config.DoubleHeaders, lgSchedule[homeTeam.ID], seriesLength)
			atNextGame := awayTeam.NextPlayableDate(config.StartDate, config.DoubleHeaders, lgSchedule[awayTeam.ID], seriesLength)
			nextGame := maxTime(htNextGame, atNextGame)

			newGame := game{ID: lgGameID, AwayTeam: homeTeam.ID, HomeTeam: awayTeam.ID, Time: nextGame}
			lgSchedule[homeTeam.ID] = append(lgSchedule[homeTeam.ID], newGame)
			lgSchedule[awayTeam.ID] = append(lgSchedule[awayTeam.ID], newGame)
			lgGames = append(lgGames, newGame)
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
