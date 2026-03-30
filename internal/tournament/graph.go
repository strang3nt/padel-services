package tournament

import (
	"iter"
	"maps"
)

type Graph struct {
	nodes map[Node]map[Node]bool
}

type Node int
type edge struct {
	P1 Node
	P2 Node
}

type matching map[edge]struct{}
type matchings []matching

func (m *matching) addCanonicalEdge(nodeA int, nodeB int) {
	if *m == nil {
		*m = make(matching)
	}

	var e edge
	if nodeA < nodeB {
		e = edge{P1: Node(nodeA), P2: Node(nodeB)}
	} else {
		e = edge{P1: Node(nodeB), P2: Node(nodeA)}
	}

	(*m)[e] = struct{}{}
}

func copyMatchings(src matchings) matchings {
	dst := make(matchings, len(src))
	for i, m := range src {
		dst[i] = make(matching)
		maps.Copy(dst[i], m)
	}
	return dst
}

func kRegularEven(nodes []int, k int) matching {
	n := len(nodes)
	res := make(matching)

	for i := range n {

		for count := 1; count <= k/2; count++ {
			jIndex := i - count

			jModN := (jIndex%n + n) % n

			res.addCanonicalEdge(nodes[jModN], nodes[i])
		}

		for count := 1; count <= k/2; count++ {
			jIndex := i + count
			jModN := jIndex % n

			res.addCanonicalEdge(nodes[jModN], nodes[i])
		}
	}

	return res
}

func kRegularOdd(nodes []int, k int) matching {
	n := len(nodes)

	res := kRegularEven(nodes, k-1)

	for i := range n {

		partnerIndex := (i + n/2) % n
		res.addCanonicalEdge(nodes[i], nodes[partnerIndex])
	}

	return res
}

func makeMatching(nodes []int, k int) matching {
	n := len(nodes)

	isEvenNK := (n*k)%2 == 0

	isNGreaterThanK := n > k

	if !isEvenNK || !isNGreaterThanK {
		return make(matching)
	}

	if k%2 == 0 {

		return kRegularEven(nodes, k)
	}

	return kRegularOdd(nodes, k)
}

type nodeSet map[int]any

func (ns nodeSet) contains(node int) bool {
	_, ok := ns[node]
	return ok
}
func (g Graph) Empty() bool {

	for _, neighbors := range g.nodes {
		if len(neighbors) > 0 {
			return false
		}
	}

	return true
}

func (g *Graph) RemoveEdge(e edge) bool {

	res := false

	if neighbors, exists := g.nodes[Node(e.P1)]; exists {
		delete(neighbors, Node(e.P2))
		res = true
	}
	if neighbors, exists := g.nodes[Node(e.P2)]; exists {
		delete(neighbors, Node(e.P1))
		res = true
	}

	return res
}

func (g Graph) GetCopy() Graph {
	newGraph := Graph{
		nodes: make(map[Node]map[Node]bool),
	}

	for node, neighbors := range g.nodes {
		newNeighbors := make(map[Node]bool)
		for neighbor := range neighbors {
			newNeighbors[neighbor] = true
		}
		newGraph.nodes[node] = newNeighbors
	}

	return newGraph
}

func MakeGraph() Graph {
	return Graph{
		nodes: make(map[Node]map[Node]bool),
	}
}

func (g *Graph) GetNeighbors(n Node) []Node {
	neighborsList := []Node{}
	if neighbors, exists := g.nodes[n]; exists {
		for neighbor := range neighbors {
			neighborsList = append(neighborsList, neighbor)
		}
	}
	return neighborsList
}

func (g Graph) GetAdjacentEdges(n Node) []edge {
	var edges []edge
	if neighbors, exists := g.nodes[n]; exists {
		for neighbor := range neighbors {
			edges = append(edges, edge{P1: n, P2: neighbor})
		}
	}
	return edges
}

func (g Graph) Size() int {
	count := 0
	for _, neighbors := range g.nodes {
		count += len(neighbors)
	}
	return count / 2
}

func (g Graph) GetEdgesIterator() iter.Seq[edge] {

	return func(yield func(edge) bool) {
		seen := make(map[edge]bool)

		for node, neighbors := range g.nodes {
			for neighbor := range neighbors {
				e := edge{P1: node, P2: neighbor}
				eRev := edge{P1: neighbor, P2: node}
				if !seen[e] && !seen[eRev] {
					if !yield(e) {
						return
					}
					seen[e] = true
				}
			}
		}
	}
}

func (g *Graph) AddEdge(e edge) {
	if _, exists := g.nodes[e.P1]; !exists {
		g.nodes[e.P1] = make(map[Node]bool)
	}
	g.nodes[e.P1][e.P2] = true

	if _, exists := g.nodes[e.P2]; !exists {
		g.nodes[e.P2] = make(map[Node]bool)
	}
	g.nodes[e.P2][e.P1] = true
}
