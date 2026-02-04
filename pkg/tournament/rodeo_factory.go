package tournament

import (
	"context"
	"errors"
	"fmt"
	"log"
	"maps"
	"math"
	"time"
)

type RodeoFactory struct {
	MaxRounds       int
	AvailableCourts int
}

func (rf *RodeoFactory) GetFirstValidTournament(
	timeout time.Duration,
	count int,
	teams []Team,
	start time.Time,
) (*Rodeo, error) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan *Rodeo, count)

	for i := range count {
		go func(id int) {

			rodeo, err := rf.MakeTournament(ctx, teams, start)
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

func (rf *RodeoFactory) MakeTournament(
	ctx context.Context,
	teams []Team,
	dateStart time.Time) (*Rodeo, error) {
	n := len(teams)
	nodes := make([]int, n)
	for i := range n {
		nodes[i] = i
	}

	totalMatches, matchesPerTurn, matchesPerTeam :=
		getMatchesPerTeam(n, rf.MaxRounds, rf.AvailableCourts)

	roundsNumber := rf.MaxRounds

	if math.Ceil(matchesPerTurn)*float64(rf.MaxRounds)-1 > float64(totalMatches) {
		roundsNumber = rf.MaxRounds - 1
	}

	if totalMatches == 0 {
		return nil, errors.New("could not determine valid match parameters. Returning empty tournament")
	}

	graph, teams := rf.getGraph(teams, matchesPerTeam)

	var rounds matchings
	var err error
	rounds = rf.makeMatchingsHeuristic(*graph.GetCopy(), matchesPerTurn, roundsNumber)

	noRoundsAreEmpty := true
	for _, round := range rounds {
		if len(round) == 0 {
			noRoundsAreEmpty = false
			break
		}
	}

	if !noRoundsAreEmpty || len(rounds) != roundsNumber {
		rounds, err = rf.makeMatchingsBacktracking(ctx, *graph, matchesPerTurn, roundsNumber)
		if err != nil {
			return nil, err
		}
	}

	err = validateTournamentRounds(rounds, teams, roundsNumber, matchesPerTurn, totalMatches, matchesPerTeam)
	if err != nil {
		log.Printf("Validation error: %v", err)
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
				TeamA:   teamA,
				TeamB:   teamB,
				CourtId: currCourt,
			}
			currCourt += 1
			matches = append(matches, m)
		}

		turns = append(turns, Round{matches})
	}

	return NewRodeo("Rodeo", dateStart, teams, turns), nil
}

func (rf *RodeoFactory) getGraph(teams []Team, matchesPerTeam int) (*Graph, []Team) {

	graph := NewGraph()
	allMatches := make(matching)
	var teamsOrdered []Team

	if canAllGendersPlayOnlyAgainstEachOther(teams, matchesPerTeam) {
		teamsSeparatedByGender := make([]Team, 0, len(teams))
		nodesByGender := make(map[Gender][]int)
		currNode := 0
		for _, g := range GetAllGenders() {
			genderTeams := GetTeamsByGender(teams, g)
			teamsSeparatedByGender = append(teamsSeparatedByGender, genderTeams...)
			nodesByGender[g] = make([]int, 0, len(genderTeams))
			for i := currNode; i < currNode+len(genderTeams); i++ {
				nodesByGender[g] = append(nodesByGender[g], i)
			}
			currNode += len(genderTeams)
		}

		for _, v := range nodesByGender {
			for edge := range rf.makeEdges(v, matchesPerTeam) {
				allMatches[edge] = struct{}{}
			}
		}
		teamsOrdered = teamsSeparatedByGender

	} else {
		nodes := make([]int, len(teams))
		for i := range len(teams) {
			nodes[i] = i
		}
		allMatches = rf.makeEdges(nodes, matchesPerTeam)
		teamsOrdered = orderTeamsByGender(teams)
	}

	for edge := range allMatches {
		graph.AddEdge(edge)
	}
	return graph, teamsOrdered
}

func canAllGendersPlayOnlyAgainstEachOther(teams []Team, matchesPerTeam int) bool {
	genderCounts := make(map[Gender]int)

	for _, g := range GetAllGenders() {
		genderCounts[g] = 0
	}

	for _, team := range teams {
		genderCounts[team.TeamGender] += 1
	}

	allGendersCanPlayOnlyAgainstEachOther := true
	for _, n := range genderCounts {
		if n <= matchesPerTeam {
			allGendersCanPlayOnlyAgainstEachOther = false
		}
	}

	return allGendersCanPlayOnlyAgainstEachOther
}

func getMatchesPerTeam(teamsNumber int, totalRounds int, availableCourts int) (int, float64, int) {

	matchesPerTeam := totalRounds

	for matchesPerTeam > 0 {
		totalParticipations := teamsNumber * matchesPerTeam

		totalMatchesFloat := float64(totalParticipations) / 2.0

		if math.Floor(totalMatchesFloat) == totalMatchesFloat {

			matchesPerTurn := totalMatchesFloat / float64(totalRounds)

			if matchesPerTurn <= float64(availableCourts) && teamsNumber > matchesPerTeam {
				return int(totalMatchesFloat), matchesPerTurn, matchesPerTeam
			}
		}

		matchesPerTeam -= 1
	}

	return 0, 0.0, 0
}

type Node int
type edge struct {
	P1 Node
	P2 Node
}

type matching map[edge]struct{}
type matchings []matching

type nodeSet map[int]struct{}

func (ns nodeSet) contains(node int) bool {
	_, ok := ns[node]
	return ok
}

func (rf *RodeoFactory) makeMatchingsHeuristic(
	graph Graph,
	avgMatchingSize float64,
	totalMatchings int,
) matchings {
	var res matchings
	maxMatchesPerTurn := int(math.Ceil(avgMatchingSize))

	playingTeams := make(nodeSet)
	turnMatches := make(matching)

	for len(res) < totalMatchings {
		addedSomething := false

		for i := range graph.nodes {
			neighbors := graph.GetNeighbors(i)

			var player2 Node
			foundPlayer2 := false

			for _, neighbor := range neighbors {
				if _, playing := playingTeams[int(neighbor)]; !playing {
					player2 = neighbor
					foundPlayer2 = true
					break
				}
			}

			if !playingTeams.contains(int(i)) && foundPlayer2 {
				addedSomething = true
				player1 := i

				playingTeams[int(player1)] = struct{}{}
				playingTeams[int(player2)] = struct{}{}

				turnMatches[edge{player1, player2}] = struct{}{}

				graph.RemoveEdge(edge{player1, Node(player2)})

			}

			isLastTurnRequirement := len(res) == totalMatchings-1 && graph.Empty()

			if len(turnMatches) == maxMatchesPerTurn || isLastTurnRequirement {
				if len(turnMatches) > 0 {
					res = append(res, turnMatches)
					playingTeams = make(nodeSet)
					turnMatches = make(matching)
				}
			}

		}

		if !addedSomething {
			break
		}
	}

	return res
}

func (rf *RodeoFactory) addCanonicalEdge(res matching, nodeA int, nodeB int) {
	var e edge
	if nodeA < nodeB {
		e = edge{P1: Node(nodeA), P2: Node(nodeB)}
	} else {
		e = edge{P1: Node(nodeB), P2: Node(nodeA)}
	}
	res[e] = struct{}{}
}

func (rf *RodeoFactory) kRegularEven(nodes []int, k int) matching {
	n := len(nodes)
	res := make(matching)

	for i := range n {

		for count := 1; count <= k/2; count++ {
			jIndex := i - count

			jModN := (jIndex%n + n) % n

			rf.addCanonicalEdge(res, nodes[jModN], nodes[i])
		}

		for count := 1; count <= k/2; count++ {
			jIndex := i + count
			jModN := jIndex % n

			rf.addCanonicalEdge(res, nodes[jModN], nodes[i])
		}
	}

	return res
}

func (rf *RodeoFactory) kRegularOdd(nodes []int, k int) matching {
	n := len(nodes)

	res := rf.kRegularEven(nodes, k-1)

	for i := range n {

		partnerIndex := (i + n/2) % n
		rf.addCanonicalEdge(res, nodes[i], nodes[partnerIndex])
	}

	return res
}

func (rf *RodeoFactory) makeEdges(nodes []int, k int) matching {
	n := len(nodes)

	isEvenNK := (n*k)%2 == 0

	isNGreaterThanK := n > k

	if !isEvenNK || !isNGreaterThanK {
		return make(matching)
	}

	if k%2 == 0 {

		return rf.kRegularEven(nodes, k)
	}

	return rf.kRegularOdd(nodes, k)
}

func NewRodeoFactory(turns, availableCourts int) *RodeoFactory {
	return &RodeoFactory{
		MaxRounds:       turns,
		AvailableCourts: availableCourts,
	}
}

func prepend(tms []Team, t Team) []Team {
	tms = append(tms, Team{})
	copy(tms[1:], tms)
	tms[0] = t
	return tms
}

// Orders the teams in such a way that Female teams are in the middle,
// surrounded by Else (mixed) teams, and Male teams are at the extremes.
// This is done to favor teams playing against same gender teams. Note that
// the ordering is built like this because it leverages some aspects of the
// the match-making algorithm: teams will be paired against the nearest neighbors
// to their left and right.
func orderTeamsByGender(teams []Team) []Team {

	genderBuckets := make(map[Gender][]Team)

	for _, team := range teams {
		genderBuckets[team.TeamGender] = append(genderBuckets[team.TeamGender], team)
	}

	var orderedTeams []Team
	orderedTeams = append(orderedTeams, genderBuckets[Female]...)

	top := true
	for _, team := range genderBuckets[Else] {
		if top {
			orderedTeams = prepend(orderedTeams, team)
			top = false
		} else {
			orderedTeams = append(orderedTeams, team)
			top = true
		}
	}

	top = true
	for _, team := range genderBuckets[Male] {
		if top {
			orderedTeams = prepend(orderedTeams, team)
			top = false
		} else {
			orderedTeams = append(orderedTeams, team)
			top = true
		}
	}

	return orderedTeams
}

func (rf *RodeoFactory) makeMatchingsBacktracking(
	ctx context.Context, initialEdges Graph, avgMatchingSize float64, totalMatchings int) (matchings, error) {

	maxMatchingSize := int(math.Ceil(avgMatchingSize))

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

	result, success := rf.solveRecursive(ctx, edgeList, remainingEdges, 0, buckets, usedNodes, maxMatchingSize)

	if success {
		return result, nil
	}

	return nil, errors.New("could not find valid matchings with the given parameters")
}

func (rf *RodeoFactory) solveRecursive(
	ctx context.Context,
	allEdges []edge,
	remainingEdges map[edge]struct{},
	edgeIdx int,
	buckets matchings,
	usedNodes []nodeSet,
	maxSize int) (matchings, bool) {

	select {
	case <-ctx.Done():
		return nil, false
	default:
	}

	if edgeIdx == len(allEdges) {
		return copyMatchings(buckets), true
	}

	currentEdge := allEdges[edgeIdx]
	minOptions := len(buckets) + 1

	// Purpose of the loop: fail quickly. If any edge has no options, backtrack
	// immediately. Otherwise choose the edge that is the most difficult to place,
	// in an effor to reduce the branching factor of the search tree.
	for e := range remainingEdges {
		options := 0
		for i := range buckets {
			if !usedNodes[i].contains(int(e.P1)) && !usedNodes[i].contains(int(e.P2)) && len(buckets[i]) < maxSize {
				options++
			}
		}

		if options == 0 {
			return nil, false
		}

		if options < minOptions {
			minOptions = options
			currentEdge = e
		}
		if minOptions == 1 {
			break
		}
	}

	p1 := int(currentEdge.P1)
	p2 := int(currentEdge.P2)

	for i := range buckets {

		bucketEdges := buckets[i]
		nodesInBucket := usedNodes[i]

		if !nodesInBucket.contains(p1) &&
			!nodesInBucket.contains(p2) &&
			len(bucketEdges) < maxSize {

			bucketEdges[currentEdge] = struct{}{}
			nodesInBucket[p1] = struct{}{}
			nodesInBucket[p2] = struct{}{}

			delete(remainingEdges, currentEdge)
			sol, found := rf.solveRecursive(ctx, allEdges, remainingEdges, edgeIdx+1, buckets, usedNodes, maxSize)
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

func copyMatchings(src matchings) matchings {
	dst := make(matchings, len(src))
	for i, m := range src {
		dst[i] = make(matching)
		maps.Copy(dst[i], m)
	}
	return dst
}

func validateTournamentRounds(
	rounds matchings,
	teams []Team,
	totalRounds int,
	matchesPerTurn float64,
	totalMatches int,
	matchesPerTeam int) error {
	if len(rounds) != totalRounds {
		return fmt.Errorf("expected %d rounds, got %d", totalRounds, len(rounds))
	}

	maxMatchesPerTurn := int(math.Ceil(matchesPerTurn))

	for i, round := range rounds {
		if len(round) > maxMatchesPerTurn {
			return fmt.Errorf("round %d violated constraint: expected <= %d matches, got %d",
				i+1, maxMatchesPerTurn, len(round))
		}
	}

	scheduledEdges := make(matching)
	totalScheduledCount := 0

	for i, round := range rounds {
		for edge := range round {
			totalScheduledCount++
			if _, exists := scheduledEdges[edge]; exists {
				return fmt.Errorf("match scheduled twice: teams %v and %v were scheduled again in round %d",
					teams[edge.P1], teams[edge.P2], i+1)
			}
			scheduledEdges[edge] = struct{}{}
		}
	}

	if totalScheduledCount != totalMatches {
		return fmt.Errorf("not every match was scheduled: total scheduled count %d does not match total matches %d",
			totalScheduledCount, totalMatches)
	}

	scheduledNodes := make(map[Node]int)
	for i, round := range rounds {

		scheduledInRound := make(map[Node]struct{})
		for edge := range round {

			if _, exists := scheduledInRound[edge.P1]; exists {
				return fmt.Errorf("team %v scheduled more than once across all rounds (found twice in round %d)", teams[edge.P1], i+1)
			}

			if _, exists := scheduledInRound[edge.P2]; exists {
				return fmt.Errorf("team %v scheduled more than once across all rounds (found twice in round %d)", teams[edge.P2], i+1)
			}

			scheduledInRound[edge.P1] = struct{}{}
			scheduledInRound[edge.P2] = struct{}{}
			scheduledNodes[edge.P1] += 1
			scheduledNodes[edge.P2] += 1
		}
	}

	for node, count := range scheduledNodes {

		if count != matchesPerTeam {
			return fmt.Errorf("team %v scheduled %d times, expected %d times", teams[node], count, matchesPerTeam)
		}
	}

	return nil
}
