package tournament

import (
	"fmt"
	"time"
)

type Round struct {
	Matches []Match `json:"matches"`
}

type MatchStatus int

const (
	MatchScheduled MatchStatus = iota
	MatchOngoing
	MatchCompleted
)

type TournamentType int

const (
	TournamentTypeEmpty TournamentType = iota
	TournamentTypeRodeo
)

type Match struct {
	TeamA       Team        `json:"teamA"`
	TeamB       Team        `json:"teamB"`
	MatchStatus MatchStatus `json:"matchStatus"`
	CourtId     int         `json:"courtId"`
}

type Tournament interface {
	GetName() string
	GetDateStart() time.Time
	GetTeams() []Team
	GetRounds() []Round
	GetTournamentType() TournamentType
}

type TournamentData struct {
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	Teams  []Team    `json:"teams"`
	Rounds []Round   `json:"rounds"`
}

func MakeTournamentData(name string, date time.Time, teams []Team, rounds []Round) TournamentData {
	return TournamentData{
		name, date, teams, rounds,
	}
}

func TournamentTypeToString(t TournamentType) (string, error) {
	switch t {
	case TournamentTypeRodeo:
		return "Rodeo", nil
	default:
		return "", fmt.Errorf("invalid tournament type: %d", t)
	}
}

type TournamentFactory interface {
	MakeTournament(teams []Team, dateStart time.Time) (*Tournament, error)
}
