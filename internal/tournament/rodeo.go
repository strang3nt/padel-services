package tournament

import (
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
