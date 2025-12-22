package tournament

import (
	"errors"
	"math"
	"time"
)

type RodeoFactory struct {
	TotalRounds     int
	AvailableCourts int
}

func (rf *RodeoFactory) MakeTournament(
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

	rounds, err := rf.makeMatchingsBruteForce(allMatches, matchesPerTurn, rf.TotalRounds)
	if err != nil {
		return nil, err
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

			if matchesPerTurn <= float64(availableCourts) {
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

type state struct {
	buckets          matchings
	remainingEdges   matching // Set of remaining edges
	bucketsUsedNodes []nodeSet
}

type graphState struct {
	buckets          matchings
	remainingEdges   Graph // Set of remaining edges
	bucketsUsedNodes []nodeSet
}

func (ns nodeSet) contains(node int) bool {
	_, ok := ns[node]
	return ok
}

func removeEdge(m matching, e edge) bool {
	if _, exists := m[e]; exists {
		delete(m, e)
		return true
	}
	return false
}

func copyState(s *state) *state {
	newState := &state{
		buckets:          make(matchings, len(s.buckets)),
		remainingEdges:   make(matching, len(s.remainingEdges)),
		bucketsUsedNodes: make([]nodeSet, len(s.bucketsUsedNodes)),
	}

	for i, m := range s.buckets {
		newState.buckets[i] = make(matching, len(m))
		for edge := range m {
			newState.buckets[i][edge] = struct{}{}
		}
	}

	for edge := range s.remainingEdges {
		newState.remainingEdges[edge] = struct{}{}
	}

	for i, ns := range s.bucketsUsedNodes {
		newState.bucketsUsedNodes[i] = make(nodeSet, len(ns))
		for node := range ns {
			newState.bucketsUsedNodes[i][node] = struct{}{}
		}
	}

	return newState
}

func copyGraphState(s *graphState) *graphState {
	newState := &graphState{
		buckets:          make(matchings, len(s.buckets)),
		remainingEdges:   *s.remainingEdges.GetCopy(),
		bucketsUsedNodes: make([]nodeSet, len(s.bucketsUsedNodes)),
	}

	for i, m := range s.buckets {
		newState.buckets[i] = make(matching, len(m))
		for edge := range m {
			newState.buckets[i][edge] = struct{}{}
		}
	}

	for i, ns := range s.bucketsUsedNodes {
		newState.bucketsUsedNodes[i] = make(nodeSet, len(ns))
		for node := range ns {
			newState.bucketsUsedNodes[i][node] = struct{}{}
		}
	}

	return newState
}

func (rf *RodeoFactory) makeMatchingsBruteForceGraph(
	initialEdges Graph, avgMatchingSize float64, totalMatchings int) (matchings, error) {

	maxMatchingSize := int(math.Ceil(avgMatchingSize))

	var stack []*graphState

	initialBuckets := make(matchings, totalMatchings)
	for i := range initialBuckets {
		initialBuckets[i] = make(matching)
	}

	initialUsedNodes := make([]nodeSet, totalMatchings)
	for i := range initialUsedNodes {
		initialUsedNodes[i] = make(nodeSet)
	}

	initialState := &graphState{
		buckets:          initialBuckets,
		remainingEdges:   initialEdges,
		bucketsUsedNodes: initialUsedNodes,
	}

	stack = append(stack, initialState)

	for len(stack) > 0 {

		currentState := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if currentState.remainingEdges.Size() == 0 {
			return currentState.buckets, nil
		}

		var currentEdge edge
		for edge := range currentState.remainingEdges.GetEdgesIterator() {
			currentEdge = edge
			break
		}
		p1 := currentEdge.P1
		p2 := currentEdge.P2

		for i := range totalMatchings {
			players := currentState.bucketsUsedNodes[i]
			edgesInMatching := int(len(currentState.buckets[i]))

			if !players.contains(int(p1)) && !players.contains(int(p2)) &&
				edgesInMatching < maxMatchingSize {

				nextState := copyGraphState(currentState)

				nextState.buckets[i][currentEdge] = struct{}{}
				nextState.bucketsUsedNodes[i][int(p1)] = struct{}{}
				nextState.bucketsUsedNodes[i][int(p2)] = struct{}{}

				if !nextState.remainingEdges.RemoveEdge(currentEdge) {
					return nil, errors.New("edge was not found in remaining_edges during removal")
				}

				stack = append(stack, nextState)
			}
		}
	}

	return nil, errors.New("could not find valid matchings with the given parameters")
}

func (rf *RodeoFactory) makeMatchingsBruteForce(
	initialEdges matching, avgMatchingSize float64, totalMatchings int) (matchings, error) {

	maxMatchingSize := int(math.Ceil(avgMatchingSize))

	var stack []*state

	initialBuckets := make(matchings, totalMatchings)
	for i := range initialBuckets {
		initialBuckets[i] = make(matching)
	}

	initialUsedNodes := make([]nodeSet, totalMatchings)
	for i := range initialUsedNodes {
		initialUsedNodes[i] = make(nodeSet)
	}

	initialState := &state{
		buckets:          initialBuckets,
		remainingEdges:   initialEdges,
		bucketsUsedNodes: initialUsedNodes,
	}

	stack = append(stack, initialState)

	for len(stack) > 0 {

		currentState := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if len(currentState.remainingEdges) == 0 {
			return currentState.buckets, nil
		}

		var currentEdge edge
		for edge := range currentState.remainingEdges {
			currentEdge = edge
			break
		}
		p1 := currentEdge.P1
		p2 := currentEdge.P2

		for i := range totalMatchings {
			players := currentState.bucketsUsedNodes[i]
			edgesInMatching := int(len(currentState.buckets[i]))

			if !players.contains(int(p1)) && !players.contains(int(p2)) &&
				edgesInMatching < maxMatchingSize {

				nextState := copyState(currentState)

				nextState.buckets[i][currentEdge] = struct{}{}
				nextState.bucketsUsedNodes[i][int(p1)] = struct{}{}
				nextState.bucketsUsedNodes[i][int(p2)] = struct{}{}

				if !removeEdge(nextState.remainingEdges, currentEdge) {
					return nil, errors.New("edge was not found in remaining_edges during removal")
				}

				stack = append(stack, nextState)
			}
		}
	}

	return nil, errors.New("could not find valid matchings with the given parameters")
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

	if !(isEvenNK && isNGreaterThanK) {
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
