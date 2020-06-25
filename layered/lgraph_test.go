package layered

import (
	"testing"

	"github.com/vibridi/graphly"

	"github.com/stretchr/testify/assert"
	"github.com/vibridi/graphly/internal"
)

func TestToLayeredGraph(t *testing.T) {
	cases := internal.ReadTestFilesFirstn(dirFlattened, 5)

	for _, test := range cases {
		g := graphly.FromJson(test.Data, test.Name)
		adj := g.AdjacencyList()
		lg := toLayeredGraph(g)

		for i, wantN := range g.Children {
			gotN := lg.Nodes[i]
			// Nodes correspond to each other
			assert.Equal(t, wantN.ID, gotN.id)
			for j, wantP := range wantN.Ports {
				// Ports correspond to each other
				assert.Equal(t, wantP.ID, gotN.ports[j].id)

				for _, gotP := range gotN.ports {
					// Port owner is correct
					assert.Equal(t, gotN.id, gotP.owner.id)

					// Outgoing edges are correct
					for _, gotOE := range gotP.outEdges {
						assert.Equal(t, gotP.id, gotOE.source.id)
						assert.Contains(t, adj[gotN.id], gotOE.target.owner.id)
					}

					// Incoming edges are correct
					for _, gotIE := range gotP.inEdges {
						assert.Equal(t, gotP.id, gotIE.target.id)
						assert.Contains(t, adj[gotIE.source.owner.id], gotIE.target.owner.id)
					}
				}
			}
		}
	}
}

func TestSplit(t *testing.T) {
	test := internal.ReadTestFile(dirCyclic, "aspect_cartrackingattackmodeling_CarTrackingAttackModeling.json")
	lg := toLayeredGraph(graphly.FromJson(test.Data, test.Name))

	parts := split(lg)
	assert.Len(t, parts, 5)

	cases := []int{17, 8, 8, 11, 6}

	for i, p := range parts {
		assert.Len(t, p.Nodes, cases[i])
	}
}
