package graphly

// graph type

type Graph struct {
	nodes    nodelist
	isCyclic bool
}

type Node struct {
	id    int
	layer int
	ports []*Port
}

type PortSide uint8

const (
	PortSide_NORTH PortSide = iota
	PortSide_WEST
	PortSide_SOUTH
	PortSide_EAST
)

type Port struct {
	owner    *Node
	side     PortSide
	anchor   point
	inEdges  []*Edge // incoming edges
	outEdges []*Edge // outgoing edges
}

// todo name
func (p Port) Len() int {
	return len(p.inEdges) + len(p.outEdges)
}

func (p Port) edge(i int) *Edge {
	switch {
	case i < len(p.inEdges):
		return p.inEdges[i]

	case i < len(p.outEdges)+len(p.inEdges):
		return p.outEdges[i-len(p.inEdges)]

	default:
		return nil
	}
}

type Edge struct {
	bends             []point
	source            *Port
	target            *Port
	isReversed        bool
	PriorityDirection uint
}

type point struct {
	x float32
	y float32
}
