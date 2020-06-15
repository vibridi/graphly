package graphly

type Size struct {
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Meta struct {
	ID         string                 `json:"id"`
	Properties map[string]interface{} `json:"properties"`
}

type Node struct {
	Meta
	Size
	Point
	Ports    []*Port `json:"ports"`
	Children []*Node `json:"children"`
	Edges    []*Edge `json:"edges"`
}

type Port struct {
	Meta
	Size
	Point
}

type Edge struct {
	Meta
	Sources []string `json:"sources"`
	Targets []string `json:"targets"`
	// Sections
}

func (root *Node) PortMap() map[string]string {
	m := make(map[string]string, 0)
	for _, n := range root.Children {
		for _, p := range n.Ports {
			m[p.ID] = n.ID
		}
	}
	return m
}

func (root *Node) AdjacencyList() map[string][]string {
	adj := make(map[string][]string, len(root.Children))
	portmap := root.PortMap()
	for _, e := range root.Edges {
		src := portmap[e.Sources[0]]
		tgt := portmap[e.Targets[0]]
		adj[src] = append(adj[src], tgt)
	}
	return adj
}
