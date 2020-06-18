package layered

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/graphly/internal"
)

func TestGreedyCycleBreaker(t *testing.T) {
	gs := internal.ReadTestFilesAll(internal.DirCyclic)

	for _, g := range gs {
		t.Run("test "+g.ID, func(t *testing.T) {
			assert.True(t, hasCycles(g))

			lg := toLayeredGraph(g)

			for _, gpart := range split(lg) {
				assert.False(t, gpart.isCyclic)
				cb := &greedyCycleBreaker{}
				cb.process(gpart)

				reversedEdges := findReversedEdges(gpart)
				if len(reversedEdges) != 0 {
					assert.True(t, gpart.isCyclic)
				}
				for _, e := range g.Edges {
					if reversedEdges[e.ID] != nil {
						tmp := e.Sources
						e.Sources = e.Targets
						e.Targets = tmp
					}
				}
			}

			assert.False(t, hasCycles(g))
		})
	}

}

func findReversedEdges(g *Graph) map[string]*Edge {
	reversedEdges := make(map[string]*Edge, 0)
	for _, n := range g.Nodes {
		for _, p := range n.ports {
			for _, in := range p.inEdges {
				if in.isReversed {
					reversedEdges[in.id] = in
				}
			}
			for _, out := range p.outEdges {
				if out.isReversed {
					reversedEdges[out.id] = out
				}
			}
		}
	}
	return reversedEdges
}
