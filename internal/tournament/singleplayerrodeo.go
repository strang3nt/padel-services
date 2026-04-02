package tournament

import (
	"time"
)

type SinglePlayerRodeo struct {
	Name      string
	DateStart time.Time
	Teams     []Team
	Rounds    []Round
}

func (rodeo *SinglePlayerRodeo) GetName() string {
	return rodeo.Name
}

func (rodeo *SinglePlayerRodeo) GetDateStart() time.Time {
	return rodeo.DateStart
}

func (rodeo *SinglePlayerRodeo) GetTeams() []Team {
	return rodeo.Teams
}

func (rodeo *SinglePlayerRodeo) GetRounds() []Round {
	return rodeo.Rounds
}

func NewSinglePlayerRodeo(name string, dateStart time.Time, teams []Team, rounds []Round) *Rodeo {
	return &Rodeo{
		Name:      name,
		DateStart: dateStart,
		Teams:     teams,
		Rounds:    rounds,
	}
}

func MakeSinglePlayerRodeo(name string, dateStart time.Time, teams []Team, rounds []Round) Rodeo {
	return Rodeo{

		Name:      name,
		DateStart: dateStart,
		Teams:     teams,
		Rounds:    rounds,
	}
}

func (rodeo *SinglePlayerRodeo) GetTournamentType() TournamentType {
	return TournamentTypeSinglePlayerRodeo
}

func (rodeo SinglePlayerRodeo) getPeople() map[Person]any {

	people := make(map[Person]any)

	for _, t := range rodeo.Teams {
		people[t.Person1] = struct{}{}
		people[t.Person2] = struct{}{}
	}
	return make(map[Person]any)
}

func (rodeo SinglePlayerRodeo) GetResting(round int, separator string) []string {
	if round > len(rodeo.Rounds)-1 || round < 0 {
		return []string{}
	}

	people := rodeo.getPeople()

	for _, m := range rodeo.Rounds[round].Matches {
		team1person1 := m.TeamA.Person1
		team2person1 := m.TeamB.Person1
		team1person2 := m.TeamA.Person2
		team2person2 := m.TeamB.Person2

		delete(people, team1person1)
		delete(people, team2person1)
		delete(people, team1person2)
		delete(people, team2person2)
	}

	res := make([]string, 0)
	for p := range people {
		res = append(
			res,
			p.Id,
		)
	}

	return res
}
