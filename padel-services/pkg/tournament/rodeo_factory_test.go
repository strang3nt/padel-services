package tournament

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func (rf *RodeoFactory) makeEdgesN6K3() matching {

	return matching{
		{P1: 0, P2: 1}: struct{}{}, {P1: 1, P2: 2}: struct{}{}, {P1: 2, P2: 3}: struct{}{},
		{P1: 3, P2: 4}: struct{}{}, {P1: 4, P2: 5}: struct{}{}, {P1: 0, P2: 5}: struct{}{},
		{P1: 0, P2: 3}: struct{}{}, {P1: 1, P2: 4}: struct{}{}, {P1: 2, P2: 5}: struct{}{},
	}
}

func TestMakeMatchesWithGendersSplit(t *testing.T) {

	teams := []Team{
		{Person_1: Person{Id: "Team1_P1"}, Person_2: Person{Id: "Team1_P2"}, TeamGender: Male},
		{Person_1: Person{Id: "Team2_P1"}, Person_2: Person{Id: "Team2_P2"}, TeamGender: Male},
		{Person_1: Person{Id: "Team3_P1"}, Person_2: Person{Id: "Team3_P2"}, TeamGender: Male},
		{Person_1: Person{Id: "Team4_P1"}, Person_2: Person{Id: "Team4_P2"}, TeamGender: Male},
		{Person_1: Person{Id: "Team5_P1"}, Person_2: Person{Id: "Team5_P2"}, TeamGender: Male},
		{Person_1: Person{Id: "Team6_P1"}, Person_2: Person{Id: "Team6_P2"}, TeamGender: Female},
		{Person_1: Person{Id: "Team7_P1"}, Person_2: Person{Id: "Team7_P2"}, TeamGender: Female},
		{Person_1: Person{Id: "Team8_P1"}, Person_2: Person{Id: "Team8_P2"}, TeamGender: Female},
		{Person_1: Person{Id: "Team9_P1"}, Person_2: Person{Id: "Team9_P2"}, TeamGender: Female},
		{Person_1: Person{Id: "Team10_P1"}, Person_2: Person{Id: "Team10_P2"}, TeamGender: Female},
	}

	// this configuration should create 5 matches per round, 4 matches per team
	// thus triggering the gender split logic for tournament creation.
	rf := RodeoFactory{
		MaxRounds:       4,
		AvailableCourts: 5,
	}

	tournament, err := rf.GetFirstValidTournament(10*time.Second, runtime.NumCPU(), teams, time.Now())
	if err != nil {
		t.Fatalf("makeMatchesWithGendersSplit returned an error: %v", err)
	}

	t.Run("Assertion_1_TournamentIsCorrectlySplitBetweenMaleAndFemale", func(t *testing.T) {
		for _, round := range tournament.GetRounds() {
			for _, match := range round {
				teamA := teams[match.TeamA.TeamGender]
				teamB := teams[match.TeamB.TeamGender]
				if teamA.TeamGender != teamB.TeamGender {
					t.Errorf("Found mixed gender match in a gender split tournament: Team %s (%v) vs Team %s (%v)",
						teamA.Person_1.Id, teamA.TeamGender, teamB.Person_1.Id, teamB.TeamGender)
				}
			}
		}
	})

}

func TestMakeMatchingsBruteForceGraph_N6K3(t *testing.T) {

	totalRounds := 3
	matchesPerTurn := 3.0 // Max 3 courts available

	rf := &RodeoFactory{
		MaxRounds:       totalRounds,
		AvailableCourts: int(matchesPerTurn),
	}

	allMatches := rf.makeEdgesN6K3()

	graph := NewGraph()
	for edge := range allMatches {
		graph.AddEdge(edge)
	}
	ctx := context.Background()
	rounds, err := rf.makeMatchingsBacktracking(ctx, *graph, matchesPerTurn, totalRounds)
	if err != nil {
		t.Fatalf("makeMatchingsBruteForceGraph returned an error: %v", err)
	}

	t.Run("Assertion_1_TotalRounds", func(t *testing.T) {
		if len(rounds) != totalRounds {
			t.Errorf("Expected exactly %d rounds, got %d", totalRounds, len(rounds))
		}
	})

	t.Run("Assertion_2_SizeConstraint", func(t *testing.T) {
		for i, round := range rounds {
			if float64(len(round)) > matchesPerTurn {
				t.Errorf("Round %d violated constraint: Expected <= %.1f matches, got %d",
					i+1, matchesPerTurn, len(round))
			}
		}
	})

	t.Run("Assertion_3_CompletenessAndNoDuplicates", func(t *testing.T) {
		scheduledEdges := make(matching)
		totalScheduledCount := 0

		// Iterate through the results and track usage
		for i, round := range rounds {
			for edge := range round {
				totalScheduledCount++
				if _, exists := scheduledEdges[edge]; exists {
					t.Fatalf("Duplicate edge found: Edge %v (Teams %d and %d) was scheduled again in Round %d",
						edge, edge.P1, edge.P2, i+1)
				}
				scheduledEdges[edge] = struct{}{}
			}
		}

		expectedUniqueCount := len(allMatches)

		// Assertion 3a: Completeness (Did we schedule all required matches?)
		if len(scheduledEdges) != expectedUniqueCount {
			t.Errorf("Completeness check failed: Expected %d unique scheduled edges (allMatches size), got %d. %d edges were missed.",
				expectedUniqueCount, len(scheduledEdges), expectedUniqueCount-len(scheduledEdges))
		}

		// Assertion 3b: No Duplicates (Total scheduled count must exactly equal unique count)
		if totalScheduledCount != expectedUniqueCount {
			t.Errorf("Duplication check failed: Total scheduled edges (%d) did not match unique edges (%d). This should only fail if edges were missed (see 3a).",
				totalScheduledCount, expectedUniqueCount)
		}
	})
}

func TestMakeTournamentN8K8(t *testing.T) {

	teams := []Team{
		{Person_1: Person{Id: "Elena Miotto"}, Person_2: Person{Id: "Alberto Rampazzo"}, TeamGender: Male},
		{Person_1: Person{Id: "Marcos Vera"}, Person_2: Person{Id: "Santiago Alonso"}, TeamGender: Male},
		{Person_1: Person{Id: "Diego Arrieta"}, Person_2: Person{Id: "Marcelo Merino"}, TeamGender: Male},
		{Person_1: Person{Id: "Cristian Garcia"}, Person_2: Person{Id: "Jorge Torres"}, TeamGender: Male},
		{Person_1: Person{Id: "Juan Perez"}, Person_2: Person{Id: "Pedro Rodriguez"}, TeamGender: Male},
		{Person_1: Person{Id: "Maria Gomez"}, Person_2: Person{Id: "Ana Lopez"}, TeamGender: Female},
		{Person_1: Person{Id: "Laura Martinez"}, Person_2: Person{Id: "Carolina Rodriguez"}, TeamGender: Female},
		{Person_1: Person{Id: "Sofia Ramirez"}, Person_2: Person{Id: "Isabella Torres"}, TeamGender: Female},
	}

	dateStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	rodeoFactory := RodeoFactory{
		MaxRounds:       8,
		AvailableCourts: 5,
	}
	ctx := context.Background()
	rodeo, err := rodeoFactory.MakeTournament(ctx, teams, dateStart)
	if err != nil {
		t.Fatalf("makeTournament returned an error: %v", err)
	}

	t.Logf("tournament created successfully: %+v", rodeo)
	totalMatches, _, _ := getMatchesPerTeam(8, 8, 5)
	t.Run("Assertion_1_AllMatchesAssigned", func(t *testing.T) {
		count := 0
		for _, round := range rodeo.GetRounds() {
			count += len(round)
		}
		if totalMatches != count {
			t.Errorf("Expected exactly %d matches, got %d", totalMatches, count)
		}
	})

	t.Run("Assertion_2_SizeConstraint", func(t *testing.T) {
		for i, round := range rodeo.GetRounds() {
			if len(round) <= 0 {
				t.Errorf("Round %d violated constraint: Expected >= 0 matches, got %d", i+1, len(round))
			}
		}
	})

}

func TestMakeTournamentN10K8(t *testing.T) {

	teams := []Team{
		{Person_1: Person{Id: "Elena Miotto"}, Person_2: Person{Id: "Alberto Rampazzo"}, TeamGender: Male},
		{Person_1: Person{Id: "Marcos Vera"}, Person_2: Person{Id: "Santiago Alonso"}, TeamGender: Male},
		{Person_1: Person{Id: "Diego Arrieta"}, Person_2: Person{Id: "Marcelo Merino"}, TeamGender: Male},
		{Person_1: Person{Id: "Cristian Garcia"}, Person_2: Person{Id: "Jorge Torres"}, TeamGender: Male},
		{Person_1: Person{Id: "Juan Perez"}, Person_2: Person{Id: "Pedro Rodriguez"}, TeamGender: Male},
		{Person_1: Person{Id: "Maria Gomez"}, Person_2: Person{Id: "Ana Lopez"}, TeamGender: Female},
		{Person_1: Person{Id: "Laura Martinez"}, Person_2: Person{Id: "Carolina Rodriguez"}, TeamGender: Female},
		{Person_1: Person{Id: "Sofia Ramirez"}, Person_2: Person{Id: "Isabella Torres"}, TeamGender: Female},
		{Person_1: Person{Id: "Marco Gaio"}, Person_2: Person{Id: "Luigina Lodi"}, TeamGender: Female},
		{Person_1: Person{Id: "Federico Manca"}, Person_2: Person{Id: "Alberto Alberti"}, TeamGender: Female},
	}

	dateStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	rodeoFactory := RodeoFactory{
		MaxRounds:       8,
		AvailableCourts: 5,
	}
	ctx := context.Background()
	rodeo, err := rodeoFactory.MakeTournament(ctx, teams, dateStart)
	if err != nil {
		t.Fatalf("makeTournament returned an error: %v", err)
	}

	t.Logf("tournament created successfully: %+v", rodeo)

	t.Run("Assertion_1_TotalRounds", func(t *testing.T) {
		if len(rodeo.GetRounds()) != 8 {
			t.Errorf("Expected exactly 8 rounds, got %d", len(rodeo.GetRounds()))
		}
	})

	t.Run("Assertion_2_SizeConstraint", func(t *testing.T) {
		for i, round := range rodeo.GetRounds() {
			if len(round) <= 0 {
				t.Errorf("Round %d violated constraint: Expected >= 0 matches, got %d", i+1, len(round))
			}
		}
	})

}

func TestMakeTournamentMaximiseMaleMatches(t *testing.T) {

	teams := []Team{
		{Person_1: Person{Id: "Elena Miotto"}, Person_2: Person{Id: "Alberto Rampazzo"}, TeamGender: Else},
		{Person_1: Person{Id: "Marcos Vera"}, Person_2: Person{Id: "Santiago Alonso"}, TeamGender: Male},
		{Person_1: Person{Id: "Diego Arrieta"}, Person_2: Person{Id: "Marcelo Merino"}, TeamGender: Male},
		{Person_1: Person{Id: "Cristian Garcia"}, Person_2: Person{Id: "Jorgina Torres"}, TeamGender: Else},
		{Person_1: Person{Id: "Juanita Perez"}, Person_2: Person{Id: "Pedro Rodriguez"}, TeamGender: Else},
		{Person_1: Person{Id: "Maria Gomez"}, Person_2: Person{Id: "Ana Lopez"}, TeamGender: Female},
		{Person_1: Person{Id: "Laura Martinez"}, Person_2: Person{Id: "Carolina Rodriguez"}, TeamGender: Female},
		{Person_1: Person{Id: "Sofia Ramirez"}, Person_2: Person{Id: "Isabella Torres"}, TeamGender: Female},
		{Person_1: Person{Id: "Marco Gaio"}, Person_2: Person{Id: "Luigina Lodi"}, TeamGender: Else},
		{Person_1: Person{Id: "Federica Manca"}, Person_2: Person{Id: "Alberta Alberti"}, TeamGender: Female},
		{Person_1: Person{Id: "Giorgia Neri"}, Person_2: Person{Id: "Luca Bianchi"}, TeamGender: Else},
		{Person_1: Person{Id: "Francesco Russo"}, Person_2: Person{Id: "Giulio Ferrari"}, TeamGender: Male},
	}

	dateStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	rodeoFactory := RodeoFactory{
		MaxRounds:       7,
		AvailableCourts: 6,
	}
	ctx := context.Background()
	rodeo, err := rodeoFactory.MakeTournament(ctx, teams, dateStart)
	if err != nil {
		t.Fatalf("makeTournament returned an error: %v", err)
	}

	t.Logf("tournament created successfully: %+v", rodeo)

	t.Run("Assertion_2_MaximizeMaleMatches", func(t *testing.T) {
		countMalesMatches := 0
		for _, round := range rodeo.GetRounds() {
			for _, match := range round {
				teamA := match.TeamA
				teamB := match.TeamB
				if teamA.TeamGender == Male && teamB.TeamGender == Male {
					countMalesMatches++
				}
			}
		}
		if countMalesMatches != 3 {
			t.Errorf("Expected 3 male matches, got %d", countMalesMatches)
		}
	})

}
