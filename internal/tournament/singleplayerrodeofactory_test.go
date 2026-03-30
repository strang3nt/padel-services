package tournament

import (
	"runtime"
	"testing"
	"time"
)

var fourPeople = map[Person]any{
	{Id: "Tizio"}:     struct{}{},
	{Id: "Caio"}:      struct{}{},
	{Id: "Sempronio"}: struct{}{},
	{Id: "Fazio"}:     struct{}{},
}

func TestComputeMatchesPerPerson(t *testing.T) {

	inputs := []struct {
		peopleNumber    int
		totalRounds     int
		availableCourts int
	}{{4, 2, 2}, {20, 6, 4}}
	expecteds := []matchesPerPerson{
		{2, 1, 2},
		{20, 4, 4},
	}

	for i := range inputs {
		input := inputs[i]
		expected := expecteds[i]
		actual := getMatchesPerPerson(input.peopleNumber, input.totalRounds, input.availableCourts)
		if expected != actual {

			t.Fatalf("Expected %v, received %v", expected, actual)
		}
	}
}

func TestGenerateTeams(t *testing.T) {

	singlePlayerRodeoFactory := SinglePlayerRodeoFactory{
		MaxRounds:       2,
		AvailableCourts: 2,
		People:          fourPeople,
	}

	teams := singlePlayerRodeoFactory.generateTeams(2)

	t.Run("generated 4 unique teams", func(t *testing.T) {
		teamsMap := make(map[Team]any)
		for i := range teams {
			teamsMap[teams[i]] = struct{}{}
		}
		if len(teamsMap) != 4 {

			t.Errorf(
				"generated %d different teams, expected 4", len(teamsMap),
			)
		}
	})
}

func TestGenerateSinglePlayerRodeo(t *testing.T) {
	singlePlayerRodeoFactory := SinglePlayerRodeoFactory{
		MaxRounds:       2,
		AvailableCourts: 2,
		People:          fourPeople,
	}

	tournament, err := singlePlayerRodeoFactory.GetFirstValidTournament(
		10*time.Second,
		runtime.NumCPU(),
		time.Now(),
	)

	if err != nil {
		t.Fatalf("unexpected error encountered while building single player rodeo: %v", err)
	}

	t.Run("generated tournament has 4 possible teams", func(t *testing.T) {

		teams := tournament.Teams
		teamsMap := make(map[Team]any)
		for _, t := range teams {
			teamsMap[t] = struct{}{}
		}
		if len(teamsMap) != 4 {

			t.Errorf(
				"generated %d different teams, expected 4", len(teamsMap),
			)
		}
	})

	t.Run("generated tournament has 2 matches", func(t *testing.T) {

		rounds := tournament.Rounds
		totalMatches := 0
		for _, r := range rounds {
			totalMatches += len(r.Matches)
		}
		if totalMatches != 2 {

			t.Errorf(
				"generated %v matches, expected 2, with tournament %v", totalMatches, tournament,
			)
		}
	})

}
