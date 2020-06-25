package layered

import (
	"errors"
	"fmt"
	"math"

	"github.com/vibridi/graphly/internal"
)

// -----------------------
// 	Phase 1 Factory
// -----------------------

func layeredPhase1Factory(strat CycleBreakingStrategy) processor {
	switch strat {
	case CycleBreakingStrategy_GREEDY:
		return &greedyCycleBreaker{}

	case CycleBreakingStrategy_DEPTH_FIRST:
		return &depthFirstCycleBreaker{}
	}
	panic(errors.New(fmt.Sprintf("unsupported cycle breaking strategy: %d", strat)))
}

// -----------------------
// 	Greedy Cycle Breaker
// -----------------------

// Greedy cycle breaking strategy. Implements processor interface
type greedyCycleBreaker struct {
	nodes      []*Node
	indegrees  []int
	outdegrees []int
	sources    []*Node // nodes with no incoming edges
	sinks      []*Node // nodes with no outgoing edges
	arcdiag    []int
}

func (p *greedyCycleBreaker) init(nodes []*Node) {
	p.nodes = nodes
	p.indegrees = make([]int, len(nodes))
	p.outdegrees = make([]int, len(nodes))
	p.arcdiag = make([]int, len(nodes))
}

func (p *greedyCycleBreaker) isSink(id int) bool {
	return p.outdegrees[id] == 0
}

func (p *greedyCycleBreaker) addSink(n *Node) {
	p.sinks = append(p.sinks, n)
}

func (p *greedyCycleBreaker) takeSink() *Node {
	n := p.sinks[0]
	p.sinks[0] = nil
	p.sinks = p.sinks[1:]
	return n
}

func (p *greedyCycleBreaker) isSource(id int) bool {
	return p.indegrees[id] == 0
}

func (p *greedyCycleBreaker) addSource(n *Node) {
	p.sources = append(p.sources, n)
}

func (p *greedyCycleBreaker) takeSource() *Node {
	n := p.sources[0]
	p.sources[0] = nil
	p.sources = p.sources[1:]
	return n
}

// Port of ELK's (Eclipse Kernel Layout) greedy cycle breaking strategy, from:
// Peter Eades, Xuemin Lin, W. F. Smyth: A fast and effective heuristic for the feedback arc set problem
// http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.47.7745
//
// The algorithm arranges the nodes of G in an arc diagram, with source nodes to the right and sink nodes to the left.
// Then it reverses edges that point left.
func (p *greedyCycleBreaker) process(graph *Graph) {

	p.init(graph.Nodes)

	// preliminary step: compute node degrees and collect sources and sinks
	for i, node := range p.nodes {
		// the node seq is used as index for the indegrees, outdegrees, and arcdiag arrays
		node.seq = i

		for _, port := range node.ports {
			for _, inEdge := range port.inEdges {
				// ignore self-loops
				if inEdge.source.owner.id == node.id {
					continue
				}
				p.indegrees[i] += int(inEdge.PriorityDirection) + 1
			}

			for _, outEdge := range port.outEdges {
				// ignore self-loops
				if outEdge.target.owner.id == node.id {
					continue
				}
				p.outdegrees[i] += int(outEdge.PriorityDirection) + 1
			}
		}

		switch {
		case p.isSink(i):
			p.addSink(node)

		case p.isSource(i):
			p.addSource(node)
		}
	}

	// arrange nodes from left to right as in an arc diagram
	// sinks are assigned a negative value at first, which will be reversed later
	nextRight := -1
	nextLeft := 1

	unprocessed := len(p.nodes)
	for unprocessed > 0 {
		for len(p.sinks) > 0 {
			sink := p.takeSink()
			p.arcdiag[sink.seq] = nextRight
			nextRight--
			p.updateNeighbors(sink)
			unprocessed--
		}

		for len(p.sources) > 0 {
			source := p.takeSource()
			p.arcdiag[source.seq] = nextLeft
			nextLeft++
			p.updateNeighbors(source)
			unprocessed--
		}

		// find the set of unprocessed node with the largest outflow
		// the outflow is the difference between outdegree and indegree
		if unprocessed > 0 {
			maxOutflow := math.MinInt32
			maxOutSet := make([]*Node, 0)

			for _, node := range p.nodes {
				if p.arcdiag[node.seq] == 0 {
					outflow := p.outdegrees[node.seq] - p.indegrees[node.seq]
					if outflow >= maxOutflow {
						if outflow > maxOutflow {
							maxOutSet = nil
							maxOutflow = outflow
						}
						maxOutSet = append(maxOutSet, node)
					}
				}
			}
			if maxOutflow == math.MinInt32 {
				panic("graphly: max outflow was not changed") // todo error messages
			}

			// randomly select a node from the ones with maximal outflow and put it left
			maxOutN := maxOutSet[internal.RandInt(len(maxOutSet))]
			p.arcdiag[maxOutN.seq] = nextLeft
			nextLeft++
			p.updateNeighbors(maxOutN)
			unprocessed--
		}
	}

	// finally we shift sinks to the right
	shift := len(p.nodes) + 1
	for i := range p.nodes {
		if p.arcdiag[i] < 0 {
			p.arcdiag[i] += shift
		}
	}

	// now nodes are arranged from left to right based on their outflow:
	// sources (0 indegree) -- other non-sinks ordered from greatest to smallest outflow -- sinks (0 outdegree)

	// remember edges to avoid modifying slices in place
	rev := make([]*Edge, 0)

	// reverse edges that point left
	for _, node := range p.nodes {
		for _, port := range node.ports {
			for _, edge := range port.outEdges {
				if p.arcdiag[node.seq] > p.arcdiag[edge.target.owner.seq] {
					rev = append(rev, edge)
					graph.isCyclic = true
				}
			}
		}
	}

	for _, edge := range rev {
		edge.reverse()
	}
}

// Updates indegree and outdegree values of the neighbors of the given node,
// simulating its removal from the graph. the sources and sinks lists are also updated.
func (p *greedyCycleBreaker) updateNeighbors(node *Node) {
	for _, port := range node.ports {

		// simulate removal of an edge target node: the outdegree of its source nodes decreases.
		for _, inEdge := range port.inEdges {
			sourceNode := inEdge.source.owner
			if node == sourceNode {
				continue
			}

			id := sourceNode.seq
			// if the source node is still unprocessed
			if p.arcdiag[id] == 0 {
				p.outdegrees[id] -= int(inEdge.PriorityDirection) + 1
				if p.outdegrees[id] <= 0 && p.indegrees[id] > 0 {
					p.addSink(sourceNode)
				}
			}
		}

		// simulate removal of an edge source node: the indegree of its target nodes decreases.
		for _, outEdge := range port.outEdges {
			targetNode := outEdge.target.owner
			if node == targetNode {
				continue
			}

			seq := targetNode.seq
			// if the target node is still unprocessed
			if p.arcdiag[seq] == 0 {
				p.indegrees[seq] -= int(outEdge.PriorityDirection) + 1
				if p.indegrees[seq] <= 0 && p.outdegrees[seq] > 0 {
					p.sources = append(p.sources, targetNode)
				}
			}
		}
	}
}

// -----------------------
// 	DFS Cycle Breaker
// -----------------------

// Depth-first cycle breaking strategy. Implements processor interface
type depthFirstCycleBreaker struct {
	nodes   []*Node
	visited []int
	roots   []int
	sources []*Node // nodes with no incoming edges
}

func (p *depthFirstCycleBreaker) init(nodes []*Node) {
	p.nodes = nodes
	p.visited = make([]int, len(nodes))
	for i := range p.visited {
		p.visited[i] = -1
	}
	p.roots = make([]int, len(nodes))
	for i := range p.roots {
		p.roots[i] = -1
	}
}

func (p *depthFirstCycleBreaker) process(graph *Graph) {

	p.init(graph.Nodes)

	// collect all sources
	for i, node := range p.nodes {
		node.seq = i

		indeg := 0
		for _, port := range node.ports {
			indeg += len(port.inEdges)
		}
		if indeg == 0 {
			p.sources = append(p.sources, node)
		}
	}

	for _, src := range p.sources {
		p.dfs(src, 0, src.seq)
	}

	// run a dfs again on unvisited nodes
	for i := range p.visited {
		if p.visited[i] == -1 {
			// this is possible because node.seq is equal to its position in the list
			// and visited is indexed with node.seq
			n := p.nodes[i]
			p.dfs(n, 0, n.seq)
		}
	}

	// remember edges to avoid modifying slices in place
	rev := make([]*Edge, 0)

	for _, node := range p.nodes {
		for _, port := range node.ports {
			for _, edge := range port.outEdges {
				target := edge.target.owner
				if target.id == node.id {
					continue
				}
				// when the connected nodes belong to the same tree and the target was seen before the source
				if p.roots[node.seq] == p.roots[target.seq] && p.visited[target.seq] < p.visited[node.seq] {
					rev = append(rev, edge)
					graph.isCyclic = true
				}
			}
		}
	}

	for _, edge := range rev {
		edge.reverse()
	}
}

func (p *depthFirstCycleBreaker) dfs(node *Node, depth int, rootSeq int) {
	if p.visited[node.seq] != -1 {
		return
	}

	p.visited[node.seq] = depth // remember at which depth this node is seen
	p.roots[node.seq] = rootSeq // remember at which root node this dfs started

	for _, port := range node.ports {
		for _, edge := range port.outEdges {
			target := edge.target.owner
			if target.id == node.id {
				continue
			}
			p.dfs(target, depth+1, rootSeq)
		}
	}
}
