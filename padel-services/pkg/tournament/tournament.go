package tournament

import (
	"fmt"
	"time"
)

type Round []Match

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
	TeamA       Team
	TeamB       Team
	MatchStatus MatchStatus
	CourtId     int
}

type TournamentManager interface {
	ScheduleMatch(roundIndex int, match Match) error
	GetTournamentDetails() string
}

type Tournament interface {
	GetName() string
	GetDateStart() time.Time
	GetTeams() []Team
	GetRounds() []Round
	GetTournamentType() TournamentType
}

type TournamentData struct {
	Name   string
	Date   time.Time
	Teams  []Team
	Rounds []Round
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
