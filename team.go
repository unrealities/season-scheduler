package main

import "time"

type team struct {
	Conference int    `json:"conference"`
	Divison    int    `json:"division"`
	ID         int    `json:"ID"`
	Name       string `json:"name"`
	Region     string `json:"region"`
}

// nextPlayableDate prevents a team from playing more than the desired number of games on a given day
// and gives the next available date that the team can play
func (t team) nextPlayableDate(date time.Time, doubleHeaders bool, games schedule) time.Time {
	// TODO: Handle doubleHeaders
	if doubleHeaders {
		return date
	}

	// base case. If a team has no games, then their first game should be on the start of the season
	if len(games) == 0 {
		return date
	}

	// get the most recently played game (which should be the last game in the schedule slice)
	mostRecentGame := games[len(games)-1]
	mostRecentDate := mostRecentGame.Time

	// return the next available date
	// TODO: account for league holidays
	// TODO: account for too many consecutive games
	// TODO: allow for travel days
	return mostRecentDate.AddDate(0, 0, 1)
}
