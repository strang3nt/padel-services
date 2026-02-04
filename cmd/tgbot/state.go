package main

import "padelservices/pkg/tournament"

type StateMachine int

const (
	None StateMachine = iota
	Started
	TournamentCreated
	PdfRequired
)

type TournamentCreationData struct {
	TournamentType  *tournament.TournamentType
	TotalRounds     int
	AvailableCourts int
}
