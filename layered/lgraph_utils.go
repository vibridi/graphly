package layered

// Returns the adjacency list of this graph. The map is from a port to its connected ports.
func (g *Graph) adjacencyList() map[string][]string {
	adj := make(map[string][]string, len(g.Nodes))
	for _, n := range g.Nodes {
		for _, p := range n.ports {
			for _, out := range p.outEdges {
				adj[n.id] = append(adj[n.id], out.target.owner.id)
			}
		}
	}
	return adj
}

func (g *Graph) hasCycles() bool {
	adj := g.adjacencyList()
	discovered := make(map[string]bool, 0)
	finished := make(map[string]bool, 0)
	cycle := false
	for nid, _ := range adj {
		if !discovered[nid] && !finished[nid] {
			discovered, finished = visit(&cycle, adj, nid, discovered, finished)
			if cycle {
				break
			}
		}
	}
	return cycle
}

func visit(cycle *bool, adjacentNodes map[string][]string, nid string, discovered, finished map[string]bool) (map[string]bool, map[string]bool) {
	discovered[nid] = true
	for _, vid := range adjacentNodes[nid] {
		if nid == vid {
			continue // ignore self loops
		}
		if discovered[vid] {
			*cycle = true
			break
		}
		if !finished[vid] {
			discovered, finished = visit(cycle, adjacentNodes, vid, discovered, finished)
		}
	}
	delete(discovered, nid)
	finished[nid] = true

	return discovered, finished
}
