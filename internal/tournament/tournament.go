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
	TournamentTypeSinglePlayerRodeo
)

type Match struct {
	TeamA       *Team       `json:"teamA"`
	TeamB       *Team       `json:"teamB"`
	MatchStatus MatchStatus `json:"matchStatus"`
	CourtId     int         `json:"courtId"`
}

type Tournament interface {
	GetName() string
	GetDateStart() time.Time
	GetTeams() []Team
	GetRounds() []Round
	GetResting(round int, separator string) []string
	GetTournamentType() TournamentType
}

type TournamentData struct {
	Name           string         `json:"name"`
	Date           time.Time      `json:"date"`
	Teams          []Team         `json:"teams"`
	Rounds         []Round        `json:"rounds"`
	TournamentType TournamentType `json:"tournamentType"`
}

func (t TournamentData) ToTournament() Tournament {

	switch t.TournamentType {
	case TournamentTypeRodeo:
		return NewRodeo(
			t.Name,
			t.Date,
			t.Teams,
			t.Rounds,
		)
	case TournamentTypeSinglePlayerRodeo:
		return NewSinglePlayerRodeo(
			t.Name,
			t.Date,
			t.Teams,
			t.Rounds,
		)
	default:
		return nil
	}
}

func MakeTournamentData(
	name string,
	date time.Time,
	teams []Team,
	rounds []Round,
	tournamentType TournamentType,
) TournamentData {
	return TournamentData{
		name, date, teams, rounds, tournamentType,
	}
}

func TournamentTypeToString(t TournamentType) (string, error) {
	switch t {
	case TournamentTypeRodeo:
		return "Rodeo", nil
	case TournamentTypeSinglePlayerRodeo:
		return "SinglePlayerRodeo", nil
	default:
		return "", fmt.Errorf("invalid tournament type: %d", t)
	}
}

func TournamentTypeFromString(t string) (TournamentType, error) {
	switch t {
	case "Rodeo":
		return TournamentTypeRodeo, nil
	case "SinglePlayerRodeo":
		return TournamentTypeSinglePlayerRodeo, nil
	default:
		return TournamentTypeEmpty, fmt.Errorf("invalid tournament type: %s", t)
	}
}
