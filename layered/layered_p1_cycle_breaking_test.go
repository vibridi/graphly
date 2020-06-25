package layered

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/graphly/internal"
)

func TestGreedyCycleBreaker(t *testing.T) {
	cases := internal.ReadTestFilesAll(dirCyclic)

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			testCycleBreaker(t, test, &greedyCycleBreaker{})
		})
	}
}

func TestDepthFirstCycleBreaker(t *testing.T) {
	cases := internal.ReadTestFilesAll(dirCyclic)

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			testCycleBreaker(t, test, &depthFirstCycleBreaker{})
		})
	}
}

func testCycleBreaker(t *testing.T, test *internal.GraphData, proc processor) {
	lg := fromJson(test)
	assert.True(t, lg.hasCycles())

	for _, gpart := range split(lg) {
		assert.False(t, gpart.isCyclic)

		cb := &depthFirstCycleBreaker{}
		cb.process(gpart)

		reversedEdges := findReversedEdges(gpart)
		if len(reversedEdges) != 0 {
			assert.True(t, gpart.isCyclic)
		}
		assert.False(t, gpart.hasCycles())
	}
	assert.False(t, lg.hasCycles())
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
