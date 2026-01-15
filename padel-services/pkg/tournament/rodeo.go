package tournament

import (
	"fmt"
	"strings"
	"time"
)

type Rodeo struct {
	Name      string
	DateStart time.Time
	Teams     []Team
	Rounds    []Round
}

func (rodeo *Rodeo) GetName() string {
	return rodeo.Name
}

func (rodeo *Rodeo) GetDateStart() time.Time {
	return rodeo.DateStart
}

func (rodeo *Rodeo) GetTeams() []Team {
	return rodeo.Teams
}

func (rodeo *Rodeo) GetRounds() []Round {
	return rodeo.Rounds
}

func NewRodeo(name string, dateStart time.Time, teams []Team, rounds []Round) *Rodeo {
	return &Rodeo{
		Name:      name,
		DateStart: dateStart,
		Teams:     teams,
		Rounds:    rounds,
	}
}

func MakeRodeo(name string, dateStart time.Time, teams []Team, rounds []Round) Rodeo {
	return Rodeo{

		Name:      name,
		DateStart: dateStart,
		Teams:     teams,
		Rounds:    rounds,
	}
}

func (rodeo *Rodeo) SerializeToCSV() string {
	var sb strings.Builder

	turns := rodeo.Rounds

	i := 1

	for _, t := range turns {

		sb.WriteString(fmt.Sprintf("Round %d,", i))
		match := 1

		for _, m := range t {

			team1 := m.TeamA
			team2 := m.TeamB

			sb.WriteString(fmt.Sprintf("Match %d,", match))
			sb.WriteString(fmt.Sprintf("%s - %s,", team1.Person_1.Id, team1.Person_2.Id))
			sb.WriteString(fmt.Sprintf("%s - %s,", team2.Person_1.Id, team2.Person_2.Id))

			match += 1
		}

		sb.WriteString("\n")
		i += 1
	}

	csvContent := sb.String()

	return csvContent
}

func (rodeo *Rodeo) GetTournamentType() TournamentType {
	return TournamentTypeRodeo
}
