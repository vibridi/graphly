package graphly

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdjacencyList(t *testing.T) {
	root := &Node{
		Children: []*Node{
			{Meta: Meta{ID: "N1"}, Ports: []*Port{{Meta: Meta{ID: "P1"}}}},
			{Meta: Meta{ID: "N2"}, Ports: []*Port{{Meta: Meta{ID: "P2"}}}},
			{Meta: Meta{ID: "N3"}, Ports: []*Port{{Meta: Meta{ID: "P3"}}}},
			{Meta: Meta{ID: "N4"}, Ports: []*Port{{Meta: Meta{ID: "P4"}}}},
			{Meta: Meta{ID: "N5"}, Ports: []*Port{{Meta: Meta{ID: "P5"}}}},
		},
		Edges: []*Edge{
			{Sources: []string{"P1"}, Targets: []string{"P2"}},
			{Sources: []string{"P1"}, Targets: []string{"P4"}},
			{Sources: []string{"P2"}, Targets: []string{"P3"}},
			{Sources: []string{"P3"}, Targets: []string{"P5"}},
			{Sources: []string{"P1"}, Targets: []string{"P5"}},
			{Sources: []string{"P4"}, Targets: []string{"P3"}},
			{Sources: []string{"P2"}, Targets: []string{"P5"}},
		},
	}

	adj := root.AdjacencyList()
	assert.ElementsMatch(t, adj["N1"], []string{"N2", "N4", "N5"})
	assert.ElementsMatch(t, adj["N2"], []string{"N3", "N5"})
	assert.ElementsMatch(t, adj["N3"], []string{"N5"})
	assert.ElementsMatch(t, adj["N4"], []string{"N3"})
	assert.Nil(t, adj["N5"])
}

// --------------
// 	Utilities
// --------------

func hasCycles(g *Node) bool {
	adj := g.AdjacencyList()
	discovered := make(map[string]bool, 0)
	finished := make(map[string]bool, 0)
	cycle := false
	for nid, _ := range adj {
		if !discovered[nid] && !finished[nid] {
			discovered, finished = visit(&cycle, adj, nid, discovered, finished)
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
