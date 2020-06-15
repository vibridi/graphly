package layered

import (
	"github.com/vibridi/graphly"
)

// graph type

// Nodes, Ports, Labels, Edges, and Edge Sections

type Graph struct {
	Nodes    []*Node
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

// func (p Port) edge(i int) *Edge {
// 	switch {
// 	case i < len(p.inEdges):
// 		return p.inEdges[i]
//
// 	case i < len(p.outEdges)+len(p.inEdges):
// 		return p.outEdges[i-len(p.inEdges)]
//
// 	default:
// 		return nil
// 	}
// }

type Edge struct {
	bends             []point
	source            *Port
	target            *Port
	isReversed        bool
	PriorityDirection uint
}

func (e *Edge) reverse() {
	oldsrc := e.source
	oldtgt := e.target
	e.target = oldsrc
	e.source = oldtgt
	e.isReversed = !e.isReversed
}

type point struct {
	x float32
	y float32
}

func convert(g *graphly.Node) *Graph {
	return nil
}

// func FromJSON(jsonBytes []byte) (*Graph, error) {
// 	root := &jsonNode{}
// 	if err := json.Unmarshal(jsonBytes, root); err != nil {
// 		return nil, errors.Wrap(err, "graphly: failed to unmarshal json source")
// 	}
// 	return newImporter().fromJSON(root), nil
// }
//
// func newImporter() *jsonImporter {
// 	return &jsonImporter{}
// }
//
// type jsonImporter struct {
// 	nodemap map[string]*Node
// 	portmap map[string]*Port
// 	edgemap map[string]*Edge
// }
//
// // todo essentially export this and make jsonNode the standard way to declare graphs (which also supports json)
// func (j *jsonImporter) fromJSON(root *jsonNode) *Graph {
// 	// todo handle graph-level properties
//
// 	for _, n := range root.Children {
// 		j.importNode(n)
// 	}
// 	return j.build()
// }
//
// func (j *jsonImporter) importNode(jsonNode *jsonNode) {
// 	// node properties not supported
// 	// node labels not supported
//
// 	n := &Node{
// 		// add pos and size
// 		ports: j.importPorts(jsonNode.Ports),
// 	}
// 	j.nodemap[jsonNode.ID] = n
//
// 	// child nodes not supported
// }
//
// func (j *jsonImporter) importPorts(jsonPorts []*jsonPort) []*Port {
//
// }
//
// func (j *jsonImporter) importEdges(jsonEdges []*jsonEdge) {
// 	if len(jsonEdges) == 0 {
// 		return
// 	}
// 	for _, jsonEdge := range jsonEdges {
// 		if len(jsonEdge.Sources) == 0 {
// 			// todo
// 		}
// 		if len(jsonEdge.Targets) == 0 {
// 			// todo
// 		}
// 		edge := &Edge{}
//
// 		src := jsonEdge.Sources[0]
// 		if n := j.nodemap[src]; n != nil {
// 			edge.source = n
// 		}
//
// 	}
// }
//
// func (j *jsonImporter) build() *Graph {
// 	nodes := make([]*Node, 0, len(j.nodemap))
// 	for _, n := range j.nodemap {
// 		nodes = append(nodes, n)
// 	}
// 	return &Graph{
// 		Nodes: nodes,
// 	}
// }
