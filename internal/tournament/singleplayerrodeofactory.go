package tournament

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
)

type SinglePlayerRodeoFactory struct {
	MaxRounds       int
	AvailableCourts int
	People          map[Person]any
}

func (rf *SinglePlayerRodeoFactory) GetFirstValidTournament(
	timeout time.Duration,
	count int,
	start time.Time,
) (*SinglePlayerRodeo, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan *SinglePlayerRodeo, count)

	for i := range count {
		go func(id int) {

			rodeo, err := rf.MakeTournament(ctx, start)
			if err == nil && rodeo != nil {
				select {
				case resultChan <- rodeo:
				case <-ctx.Done():
				}
			}
		}(i)
	}

	select {
	case firstResult := <-resultChan:
		cancel()
		return firstResult, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("tournament generation failed or timed out: %w", ctx.Err())
	}
}

func (rf *SinglePlayerRodeoFactory) MakeTournament(
	ctx context.Context,
	dateStart time.Time) (*SinglePlayerRodeo, error) {
	n := len(rf.People)
	nodes := make([]int, n)
	for i := range n {
		nodes[i] = i
	}

	matchesPerPerson :=
		getMatchesPerPerson(n, rf.MaxRounds, rf.AvailableCourts)

	roundsNumber := rf.MaxRounds

	if matchesPerPerson.TotalMatches == 0 {
		return nil, errors.New(
			"could not determine valid match parameters. Returning empty tournament",
		)
	}

	teams := rf.generateTeams(matchesPerPerson.MatchesPerPerson)
	graph := rf.getGraph(teams)

	var rounds matchings
	var err error
	forbidden := make(map[int]map[int]any)

	for i := range teams {
		forbidden[i] = make(map[int]any)
	}

	for i := range teams {
		team1 := teams[i]
		for j := i + 1; j < len(teams); j += 1 {
			if i == j {
				continue
			}
			team2 := teams[j]

			if teamsContainSamePerson(team1, team2) {
				forbidden[i][j] = struct{}{}
				forbidden[j][i] = struct{}{}
			}
		}
	}

	rounds, err = rf.makeMatchingsBacktracking(
		ctx,
		graph,
		matchesPerPerson.MatchesPerRound,
		roundsNumber,
		forbidden,
		matchesPerPerson.TotalMatches,
	)
	if err != nil {
		return nil, err
	}

	var turns []Round
	for _, matching := range rounds {
		var matches []Match

		currCourt := 1
		for edge := range matching {
			e1 := edge.P1
			e2 := edge.P2

			teamA := teams[e1]
			teamB := teams[e2]

			m := Match{
				TeamA:   &teamA,
				TeamB:   &teamB,
				CourtId: currCourt,
			}
			currCourt += 1
			matches = append(matches, m)
		}

		turns = append(turns, Round{matches})
	}

	singlePlayerRodeo := MakeSinglePlayerRodeo(
		"Single Player Rodeo",
		dateStart,
		teams,
		turns,
	)

	return &singlePlayerRodeo, nil
}

func (rf *SinglePlayerRodeoFactory) generateTeams(matchesPerPerson int) []Team {
	teams := make([]Team, 0)

	people := make([]Person, 0, len(rf.People))
	for p := range rf.People {
		people = append(people, p)
	}

	nodes := make([]int, len(people))
	for i := range people {
		nodes[i] = i
	}

	m := makeMatching(nodes, matchesPerPerson)
	for e := range m {
		x := e.P1
		y := e.P2

		team := MakeTeam(people[x], people[y], Male)
		teams = append(teams, team)
	}
	return teams
}

func teamsContainSamePerson(l Team, r Team) bool {
	counter := make(map[Person]any)

	for _, p := range []Person{l.Person1, l.Person2, r.Person1, r.Person2} {
		counter[p] = struct{}{}
	}

	return len(counter) != 4
}

func (rf *SinglePlayerRodeoFactory) getGraph(teams []Team) Graph {

	graph := MakeGraph()
	n := len(teams)

	for i := range teams {
		for j := i + 1; j < n; j += 1 {
			if !teamsContainSamePerson(teams[i], teams[j]) {
				graph.AddEdge(edge{
					P1: Node(i),
					P2: Node(j),
				})
			}
		}
	}

	return graph
}

type matchesPerPerson struct {
	TotalMatches     int
	MatchesPerRound  int
	MatchesPerPerson int
}

func getMatchesPerPerson(
	peopleNumber int,
	totalRounds int,
	availableCourts int,
) matchesPerPerson {

	totalSlots := 4 * totalRounds * availableCourts
	k := min(totalSlots/peopleNumber, totalRounds)

	for (peopleNumber*k)%4 != 0 && k > 0 {
		k -= 1
	}

	totalMatches := (k * peopleNumber) / 4
	matchesPerRound := int(math.Ceil(float64(totalMatches) / float64(totalRounds)))

	return matchesPerPerson{
		TotalMatches:     totalMatches,
		MatchesPerRound:  matchesPerRound,
		MatchesPerPerson: k,
	}
}

func NewSinglePlayerRodeoRodeoFactory(
	turns int, participants map[Person]any,
	availableCourts int,
) *SinglePlayerRodeoFactory {
	return &SinglePlayerRodeoFactory{
		MaxRounds:       turns,
		AvailableCourts: availableCourts,
		People:          participants,
	}
}

func (rf *SinglePlayerRodeoFactory) makeMatchingsBacktracking(
	ctx context.Context,
	initialEdges Graph,
	maxMatchingSize int,
	totalMatchings int,
	forbidden map[int]map[int]any,
	targetMatches int,
) (matchings, error) {

	var edgeList []edge
	for e := range initialEdges.GetEdgesIterator() {
		edgeList = append(edgeList, e)
	}

	buckets := make(matchings, totalMatchings)
	for i := range buckets {
		buckets[i] = make(matching)
	}

	usedNodes := make([]nodeSet, totalMatchings)
	for i := range usedNodes {
		usedNodes[i] = make(nodeSet)
	}

	remainingEdges := make(map[edge]struct{})
	for _, e := range edgeList {
		remainingEdges[e] = struct{}{}
	}

	result, success := rf.solveRecursive(
		ctx,
		remainingEdges,
		0,
		buckets,
		usedNodes,
		forbidden,
		maxMatchingSize,
		targetMatches,
	)

	if success {
		return result, nil
	}

	return nil, errors.New("could not find valid matchings with the given parameters")
}

func mapIntersect(a map[int]any, b map[int]any) map[int]any {
	res := make(map[int]any)

	for n := range a {
		if _, ok := b[n]; ok {
			res[n] = struct{}{}
		}
	}

	return res
}
func (rf *SinglePlayerRodeoFactory) solveRecursive(
	ctx context.Context,
	remainingEdges map[edge]struct{},
	placedMatches int,
	buckets matchings,
	usedNodes []nodeSet,
	forbidden map[int]map[int]any,
	maxSize int,
	targetMatches int,
) (matchings, bool) {

	select {
	case <-ctx.Done():
		return nil, false
	default:
	}

	if placedMatches == targetMatches {
		return copyMatchings(buckets), true
	}

	var currentEdge edge
	minOptions := len(buckets) + 1
	foundEdge := false

	for e := range remainingEdges {
		options := 0
		for i := range buckets {
			if !usedNodes[i].contains(int(e.P1)) && !usedNodes[i].contains(int(e.P2)) &&
				len(buckets[i]) < maxSize {

				p1 := int(e.P1)
				p2 := int(e.P2)
				nodesInBucket := usedNodes[i]
				if len(mapIntersect(forbidden[p1], nodesInBucket)) == 0 &&
					len(mapIntersect(forbidden[p2], nodesInBucket)) == 0 {
					options++
				}
			}
		}

		if options > 0 && options < minOptions {
			minOptions = options
			currentEdge = e
			foundEdge = true
		}
		if minOptions == 1 {
			break
		}
	}

	if !foundEdge {
		return nil, false
	}

	p1 := int(currentEdge.P1)
	p2 := int(currentEdge.P2)

	for i := range buckets {
		bucketEdges := buckets[i]
		nodesInBucket := usedNodes[i]

		if !nodesInBucket.contains(p1) &&
			!nodesInBucket.contains(p2) &&
			len(bucketEdges) < maxSize &&
			len(mapIntersect(forbidden[p1], nodesInBucket)) == 0 &&
			len(mapIntersect(forbidden[p2], nodesInBucket)) == 0 {

			bucketEdges[currentEdge] = struct{}{}
			nodesInBucket[p1] = struct{}{}
			nodesInBucket[p2] = struct{}{}

			delete(remainingEdges, currentEdge)

			sol, found := rf.solveRecursive(
				ctx,
				remainingEdges,
				placedMatches+1,
				buckets,
				usedNodes,
				forbidden,
				maxSize,
				targetMatches,
			)
			if found {
				return sol, true
			}

			remainingEdges[currentEdge] = struct{}{}
			delete(bucketEdges, currentEdge)
			delete(nodesInBucket, p1)
			delete(nodesInBucket, p2)
		}
	}

	return nil, false
}
