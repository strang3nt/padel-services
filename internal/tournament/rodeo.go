package tournament

import (
	"fmt"
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

func (rodeo *Rodeo) GetTournamentType() TournamentType {
	return TournamentTypeRodeo
}

func (rodeo Rodeo) GetResting(round int, separator string) []string {
	if round > len(rodeo.Rounds)-1 || round < 0 {
		return []string{}
	}

	teams := make(map[Team]any)

	for _, t := range rodeo.Teams {
		teams[t] = struct{}{}
	}

	for _, m := range rodeo.Rounds[round].Matches {
		team1 := m.TeamA
		team2 := m.TeamB

		delete(teams, *team1)
		delete(teams, *team2)
	}

	res := make([]string, 0)
	for t := range teams {
		res = append(
			res,
			fmt.Sprintf("%s %s %s", t.Person1, separator, t.Person2),
		)
	}

	return res
}
