package tournament

import (
	"testing"
)

func (rf *RodeoFactory) makeEdgesN6K3() matching {

	return matching{
		{P1: 0, P2: 1}: struct{}{}, {P1: 1, P2: 2}: struct{}{}, {P1: 2, P2: 3}: struct{}{},
		{P1: 3, P2: 4}: struct{}{}, {P1: 4, P2: 5}: struct{}{}, {P1: 0, P2: 5}: struct{}{},
		{P1: 0, P2: 3}: struct{}{}, {P1: 1, P2: 4}: struct{}{}, {P1: 2, P2: 5}: struct{}{},
	}
}

func TestMakeMatchingsBruteForce_N6K3(t *testing.T) {

	totalRounds := 3
	matchesPerTurn := 3.0 // Max 3 courts available

	rf := &RodeoFactory{
		TotalRounds:     totalRounds,
		AvailableCourts: int(matchesPerTurn),
	}

	allMatches := rf.makeEdgesN6K3()

	rounds, err := rf.makeMatchingsBruteForce(allMatches, matchesPerTurn, totalRounds)
	if err != nil {
		t.Fatalf("makeMatchingsBruteForce returned an error: %v", err)
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

func TestMakeMatchingsBruteForceGraph_N6K3(t *testing.T) {

	totalRounds := 3
	matchesPerTurn := 3.0 // Max 3 courts available

	rf := &RodeoFactory{
		TotalRounds:     totalRounds,
		AvailableCourts: int(matchesPerTurn),
	}

	allMatches := rf.makeEdgesN6K3()

	graph := NewGraph()
	for edge := range allMatches {
		graph.AddEdge(edge)
	}

	rounds, err := rf.makeMatchingsBruteForceGraph(*graph, matchesPerTurn, totalRounds)
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
