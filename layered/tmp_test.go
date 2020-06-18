package layered

import (
	"fmt"
	"testing"

	"github.com/vibridi/graphly"
	"github.com/vibridi/graphly/internal"
)

func TestTemp(t *testing.T) {
	gs := internal.ReadTestFilesAll(internal.DirFlattened)
	fmt.Println("tot files", len(gs))
	for _, g := range gs {
		if hasCycles(g) {
			fmt.Println(g.ID)
		}
	}
}

func hasCycles(g *graphly.Node) bool {
	discovered := make(map[string]bool, 0)
	finished := make(map[string]bool, 0)
	adj := g.AdjacencyList()
	cycle := false
	for nid, _ := range adj {
		if !discovered[nid] && !finished[nid] {
			discovered, finished = visit(&cycle, adj, nid, discovered, finished)
		}

	}
	return cycle
}

func visit(cycle *bool, adj map[string][]string, nid string, discovered, finished map[string]bool) (map[string]bool, map[string]bool) {

	discovered[nid] = true

	for _, vid := range adj[nid] {
		if nid == vid {
			continue // ignore self loops
		}
		if discovered[vid] {
			*cycle = true
			break
		}
		if !finished[vid] {
			discovered, finished = visit(cycle, adj, vid, discovered, finished)
		}
	}

	delete(discovered, nid)
	finished[nid] = true

	return discovered, finished
}
