package services

import (
	"bufio"
	"strings"
	"testing"
	"time"
)

func TestMakeTeamsFromMessageLargeMessage(t *testing.T) {
	msg :=
		`5
8
Elena Miotto, Alberto Rampazzo
Martina Sorgato, Francesco Pariotti
Matteo Sorgato, Riccardo Sacchetto
Silvia Nevola, Gennaro Nevola
Luongo Giovanni, Donato Pellegrino
Bilora Alessandra, Ferronato Debora
Giacomo Bernardo, Tommaso Ongari
Pozzato Andrea, Zampieri Paolo
Schiesaro Giacomo, Nicoletto Federico
Chiarelli Tommaso, Ercolani Francesco`

	scanner := bufio.NewScanner(strings.NewReader(msg))
	teams, err := MakeTeamsFromMessage(scanner)

	t.Logf("teams created successfully: %+v", teams)

	t.Run("Assertion_1_NoErrors", func(t *testing.T) {
		if err != nil {
			t.Errorf("encountered error while creating tournament: %s", err.Error())
		}
	})

	t.Run("Assertion_1_CorrectNumberOfTeams", func(t *testing.T) {
		if len(teams) != 10 {
			t.Errorf("not all teams where parsed, only: %v", len(teams))
		}
	})

}

func TestCreateLargeTournament(t *testing.T) {
	msg :=
		`5
8
Elena Miotto, Alberto Rampazzo
Martina Sorgato, Francesco Pariotti
Matteo Sorgato, Riccardo Sacchetto
Silvia Nevola, Gennaro Nevola
Luongo Giovanni, Donato Pellegrino
Bilora Alessandra, Ferronato Debora
Giacomo Bernardo, Tommaso Ongari
Pozzato Andrea, Zampieri Paolo
Schiesaro Giacomo, Nicoletto Federico
Chiarelli Tommaso, Ercolani Francesco`

	scanner := bufio.NewScanner(strings.NewReader(msg))
	teams, err := MakeTeamsFromMessage(scanner)
	rodeo := CreateTournament("Rodeo", time.Now(), teams, 8, 5)

	t.Logf("tournament created successfully: %+v", rodeo)

	t.Run("Assertion_1_tournamentNotNil", func(t *testing.T) {
		if err != nil {
			t.Errorf("encountered error while creating tournament: %s", err.Error())
		}
	})

	t.Run("Assertion_1_CorrectNumberOfRounds", func(t *testing.T) {
		if len(rodeo.GetRounds()) != 8 {
			t.Errorf("must be 8 rounds, got different number: %v", len(rodeo.GetRounds()))
		}
	})
}
