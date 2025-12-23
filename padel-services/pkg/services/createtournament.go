package services

import (
	"bufio"
	"padelservices/pkg/tournament"
	"strings"
	"time"
)

func CreateTournament(
	tournamentType string,
	dateStart time.Time,
	teams []tournament.Team,
	totalRounds, availableCourts int) tournament.Tournament {

	switch tournamentType {
	case "Rodeo":
		rodeo_factory := tournament.RodeoFactory{
			TotalRounds:     totalRounds,
			AvailableCourts: availableCourts,
		}
		rodeoInstance, err := rodeo_factory.MakeTournament(teams, dateStart)
		if err != nil {
			return nil
		}

		return rodeoInstance
	default:
		return nil
	}
}

func MakeTeamsFromMessage(stringteams *bufio.Scanner) ([]tournament.Team, error) {
	var teams []tournament.Team

	for stringteams.Scan() {
		line := stringteams.Text()

		row := strings.Split(line, ",")

		if len(row) < 2 {
			continue
		}

		person1 := tournament.Person{Id: strings.TrimSpace(row[0])}
		person2 := tournament.Person{Id: strings.TrimSpace(row[1])}

		teams = append(teams, tournament.MakeTeam(person1, person2, tournament.Male))
	}

	return teams, nil
}
