package main

import (
	"fmt"
	"time"
)

type game struct {
	AwayTeam int       `json:"awayTeam"`
	HomeTeam int       `json:"homeTeam"`
	ID       int       `json:"ID"`
	Time     time.Time `json:"time"`
}

type schedule []game
type schedules []schedule

// prettyPrint outputs a human readable version of a game
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
