package tournament

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"math"
	"time"
)

type RodeoFactory struct {
	TotalRounds     int
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
	ctx context.Context, // Added context
	teams []Team,
	dateStart time.Time) (*Rodeo, error) {
	n := len(teams)
	teams = orderTeamsByGender(teams)
	nodes := make([]int, n)
	for i := range n {
		nodes[i] = i
	}

	totalMatches, matchesPerTurn, matchesPerTeam :=
		rf.getMatchesPerTeam(int(n), rf.TotalRounds, rf.AvailableCourts)

	if totalMatches == 0 {
		return nil, errors.New("could not determine valid match parameters. Returning empty tournament")
	}

	allMatches := rf.makeEdges(nodes, matchesPerTeam)
	graph := NewGraph()
	for edge := range allMatches {
		graph.AddEdge(edge)
	}

	var rounds matchings
	var err error
	rounds = rf.makeMatchingsHeuristic(*graph, matchesPerTurn, rf.TotalRounds)

	if len(rounds) < rf.TotalRounds {
		rounds, err = rf.makeMatchingsBacktracking(ctx, *graph, matchesPerTurn, rf.TotalRounds)
		if err != nil {
			return nil, err
		}
	}

	var turns []Round
	for _, matching := range rounds {
		var matches []Match

		currCourt := 0
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

		turns = append(turns, matches)
	}

	return NewRodeo("Rodeo", dateStart, teams, turns), nil
}

func (rf *RodeoFactory) getMatchesPerTeam(teamsNumber int, totalRounds int, availableCourts int) (int, float64, int) {

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
	var turnMatches matching

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

			if playingTeams.contains(int(i)) && foundPlayer2 {
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
					playingTeams = make(map[int]struct{})
					turnMatches = nil
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
		TotalRounds:     turns,
		AvailableCourts: availableCourts,
	}
}

func prepend(tms []Team, t Team) []Team {
	tms = append(tms, Team{})
	copy(tms[1:], tms)
	tms[0] = t
	return tms
}

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
		}
	}

	for _, team := range genderBuckets[Male] {
		if top {
			orderedTeams = prepend(orderedTeams, team)
			top = false
		} else {
			orderedTeams = append(orderedTeams, team)
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

	result, success := rf.solveRecursive(ctx, edgeList, 0, buckets, usedNodes, maxMatchingSize)

	if success {
		return result, nil
	}

	return nil, errors.New("could not find valid matchings with the given parameters")
}

func (rf *RodeoFactory) solveRecursive(
	ctx context.Context,
	allEdges []edge,
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

			sol, found := rf.solveRecursive(ctx, allEdges, edgeIdx+1, buckets, usedNodes, maxSize)
			if found {
				return sol, true
			}

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
