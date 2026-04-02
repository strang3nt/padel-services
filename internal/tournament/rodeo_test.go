package tournament

import (
	"reflect"
	"sort"
	"testing"
)

func TestGetResting(t *testing.T) {
	p1 := Person{Id: "P1"}
	p2 := Person{Id: "P2"}
	p3 := Person{Id: "P3"}
	p4 := Person{Id: "P4"}
	p5 := Person{Id: "P5"}
	p6 := Person{Id: "P6"}

	teamA := Team{Person1: p1, Person2: p2}
	teamB := Team{Person1: p3, Person2: p4}
	teamC := Team{Person1: p5, Person2: p6}

	teams := []Team{teamA, teamB, teamC}

	rodeo := Rodeo{
		Teams: teams,
		Rounds: []Round{
			{
				Matches: []Match{
					{TeamA: &teamA, TeamB: &teamB},
				},
			},
		},
	}

	t.Run("Valid round with resting team", func(t *testing.T) {
		expected := []string{"P5 - P6"}
		result := rodeo.GetResting(0, "-")

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Round index out of bounds (too high)", func(t *testing.T) {
		result := rodeo.GetResting(1, "-")
		if len(result) != 0 {
			t.Errorf("Expected empty slice for OOB round, got %v", result)
		}
	})

	t.Run("Round index out of bounds (negative)", func(t *testing.T) {
		result := rodeo.GetResting(-1, "-")
		if len(result) != 0 {
			t.Errorf("Expected empty slice for negative round, got %v", result)
		}
	})

	t.Run("Multiple resting teams", func(t *testing.T) {
		emptyRoundRodeo := Rodeo{
			Teams:  teams,
			Rounds: []Round{{Matches: []Match{}}},
		}

		result := emptyRoundRodeo.GetResting(0, "|")
		if len(result) != 3 {
			t.Errorf("Expected 3 teams resting, got %d", len(result))
		}

		sort.Strings(result)
		expected := []string{"P1 | P2", "P3 | P4", "P5 | P6"}
		sort.Strings(expected)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})
}
