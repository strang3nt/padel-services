package tournament

import "iter"

type Graph struct {
	// nodes is a map where the key is the node ID (Node)
	// and the value is a map acting as a set of its neighbors (e.g., an adjacency list).
	// The C++ `pair.second.empty()` suggests the value is a container of neighbors.
	nodes map[Node]map[Node]bool
}

func (g *Graph) Empty() bool {

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

func (g *Graph) GetCopy() *Graph {
	newGraph := &Graph{
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

func NewGraph() *Graph {
	return &Graph{
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

func (g *Graph) GetAdjacentEdges(n Node) []edge {
	var edges []edge
	if neighbors, exists := g.nodes[n]; exists {
		for neighbor := range neighbors {
			edges = append(edges, edge{P1: n, P2: neighbor})
		}
	}
	return edges
}

func (g *Graph) Size() int {
	count := 0
	for _, neighbors := range g.nodes {
		count += len(neighbors)
	}
	return count / 2
}

func (g *Graph) GetEdgesIterator() iter.Seq[edge] {

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
