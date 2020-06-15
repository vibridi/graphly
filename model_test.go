package graphly

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdjacencyList(t *testing.T) {
	root := &Node{
		Edges: []*Edge{
			{Sources: []string{"N1"}, Targets: []string{"N2"}},
			{Sources: []string{"N1"}, Targets: []string{"N4"}},
			{Sources: []string{"N2"}, Targets: []string{"N3"}},
			{Sources: []string{"N3"}, Targets: []string{"N5"}},
			{Sources: []string{"N1"}, Targets: []string{"N5"}},
			{Sources: []string{"N4"}, Targets: []string{"N3"}},
			{Sources: []string{"N2"}, Targets: []string{"N5"}},
		},
	}

	adj := root.AdjacencyList()
	assert.ElementsMatch(t, adj["N1"], []string{"N2", "N4", "N5"})
	assert.ElementsMatch(t, adj["N2"], []string{"N3", "N5"})
	assert.ElementsMatch(t, adj["N3"], []string{"N5"})
	assert.ElementsMatch(t, adj["N4"], []string{"N3"})
	assert.Nil(t, adj["N5"])
}
