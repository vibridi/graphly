package layered

import (
	"github.com/vibridi/graphly"
	"github.com/vibridi/graphly/internal"
)

// graph type

// Nodes, Ports, Labels, Edges, and Edge Sections

type Graph struct {
	Nodes        []*Node
	isCyclic     bool
	hasSelfLoops bool
}

// placeholder func for now
func (g *Graph) copyProperties(that *Graph) {
	g.isCyclic = that.isCyclic
	g.hasSelfLoops = that.hasSelfLoops
}

type Node struct {
	id    string  // (possibly) unique identifier
	seq   int     // sequence number used to track this node in algorithms
	layer int     // layer this node belongs to
	ports []*Port // list of ports
}

func (n *Node) String() string {
	return n.id
}

func (n *Node) addPort(p *Port) {
	p.owner = n
	n.ports = append(n.ports, p)
}

type PortSide uint8

const (
	PortSide_NORTH PortSide = iota
	PortSide_WEST
	PortSide_SOUTH
	PortSide_EAST
)

type Port struct {
	id       string
	owner    *Node
	side     PortSide
	anchor   point
	inEdges  []*Edge // incoming edges
	outEdges []*Edge // outgoing edges
}

func (p *Port) addInEdge(e *Edge) {
	p.inEdges = append(p.inEdges, e)
}

func (p *Port) removeInEdge(e *Edge) {
	j := 0
	for i := 0; i < len(p.inEdges); i++ {
		if p.inEdges[i].id != e.id {
			p.inEdges[j] = p.inEdges[i]
			j++
		}
	}
	p.inEdges = p.inEdges[:j]
}

func (p *Port) addOutEdge(e *Edge) {
	p.outEdges = append(p.outEdges, e)
}

func (p *Port) removeOutEdge(e *Edge) {
	j := 0
	for i := 0; i < len(p.outEdges); i++ {
		if p.outEdges[i].id != e.id {
			p.outEdges[j] = p.outEdges[i]
			j++
		}
	}
	p.outEdges = p.outEdges[:j]
}

type Edge struct {
	id                string
	bends             []point
	source            *Port
	target            *Port
	isReversed        bool
	PriorityDirection uint
}

func (e *Edge) String() string {
	return e.source.owner.id + "->" + e.target.owner.id
}

func (e *Edge) reverse() {
	oldsrc := e.source
	oldtgt := e.target
	oldsrc.removeOutEdge(e)
	oldtgt.removeInEdge(e)
	e.target = oldsrc
	e.source = oldtgt
	e.target.addInEdge(e)
	e.source.addOutEdge(e)
	e.isReversed = !e.isReversed
}

type point struct {
	x float32
	y float32
}

func fromJson(src *internal.GraphData) *Graph {
	return toLayeredGraph(graphly.FromJson(src.Data, src.Name))
}

func toLayeredGraph(src *graphly.Node) *Graph {
	lgraph := &Graph{}
	portmap := make(map[string]*Port, 0)
	edgemap := make(map[string]*Edge, 0)

	for _, thatNode := range src.Children {
		thisNode := &Node{id: thatNode.ID}
		for _, thatPort := range thatNode.Ports {
			thisPort := &Port{
				id:     thatPort.ID,
				anchor: point{thatPort.X, thatPort.Y},
			}
			thisNode.addPort(thisPort)
			portmap[thatPort.ID] = thisPort
		}
		lgraph.Nodes = append(lgraph.Nodes, thisNode)
	}

	for _, thatEdge := range src.Edges {
		thisEdge := edgemap[thatEdge.ID]
		if thisEdge == nil {
			thisEdge = &Edge{id: thatEdge.ID}
			edgemap[thatEdge.ID] = thisEdge
		}

		sourcePort := portmap[thatEdge.Sources[0]]
		thisEdge.source = sourcePort
		sourcePort.outEdges = append(sourcePort.outEdges, thisEdge)

		targetPort := portmap[thatEdge.Targets[0]]
		thisEdge.target = targetPort
		targetPort.inEdges = append(targetPort.inEdges, thisEdge)

		if sourcePort.id == targetPort.id {
			lgraph.hasSelfLoops = true
		}
	}

	return lgraph
}

// Split the graph into connected components. In ELK this depends on the SEPARATE_CONNECTED_COMPONENTS property
// external ports and port constraints. For now we ignore that and split regardless.
func split(lgraph *Graph) []*Graph {
	for _, n := range lgraph.Nodes {
		n.seq = 0
	}
	res := make([]*Graph, 0)
	for _, n := range lgraph.Nodes {

		// find connected components
		compNs := dfs(n, nil)

		if len(compNs) != 0 {
			compG := &Graph{}
			compG.copyProperties(lgraph)
			// set EXT_PORT_CONNECTIONS if needed
			// copy padding if needed
			// remove minimum size if needed

			for _, m := range compNs {
				compG.Nodes = append(compG.Nodes, m)
				// set graph to node if needed
			}
			res = append(res, compG)
		}
	}

	return res
}

// Runs a DFS finding the connected nodes
func dfs(node *Node, connectedNodes []*Node) []*Node {
	if node.seq != 0 {
		// already visited
		return connectedNodes
	}

	// mark the node as visited
	node.seq = 1
	connectedNodes = append(connectedNodes, node)

	// check if this node is an external port dummy if needed

	for _, p := range node.ports {
		for _, in := range p.inEdges {
			connectedNodes = dfs(in.source.owner, connectedNodes)
		}
		for _, out := range p.outEdges {
			connectedNodes = dfs(out.target.owner, connectedNodes)
		}

	}
	return connectedNodes
}
