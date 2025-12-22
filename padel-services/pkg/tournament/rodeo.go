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

	i := 1 // Round counter (C++ used i=1)

	// Outer loop: Iterate through each Turn (Round)
	for _, t := range turns {
		// Start with "Round X,"
		sb.WriteString(fmt.Sprintf("Round %d,", i))

		match := 1 // Match counter (C++ used match=1)

		// Inner loop: Iterate through each Match in the Turn
		for _, m := range t {
			// In Go, since TeamA and TeamB in Match are likely *Team pointers,
			// we avoid the C++ optional check (.value()) and access them directly.
			team1 := m.TeamA
			team2 := m.TeamB

			// Match X,
			sb.WriteString(fmt.Sprintf("Match %d,", match))

			// Team 1 Format: Person1 Name - Person2 Name,
			sb.WriteString(fmt.Sprintf("%s - %s,", team1.Person_1.Id, team1.Person_2.Id))

			// Team 2 Format: Person1 Name - Person2 Name,
			sb.WriteString(fmt.Sprintf("%s - %s,", team2.Person_1.Id, team2.Person_2.Id))

			match += 1
		}

		// Equivalent of std::endl (newline)
		sb.WriteString("\n")
		i += 1
	}

	csvContent := sb.String()

	return csvContent
}

func (rodeo *Rodeo) GetTournamentType() TournamentType {
	return TournamentTypeRodeo
}
