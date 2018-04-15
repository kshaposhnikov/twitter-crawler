package graph

type Node struct {
	Name                 string
	AssociatedNodesCount int
	AssociatedNodes      []string
}

type Graph struct {
	Nodes []Node
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: []Node{},
	}
}

func (graph *Graph) Concat(newGraph Graph) *Graph{
	for _, node := range newGraph.Nodes {
		graph.AddNode(node)
	}
	return graph
}

func (graph *Graph) AddNode(node Node) *Graph {
	graph.Nodes = append(graph.Nodes, node)
	return graph
}

func (graph *Graph) ContainsVertex(name string) bool {
	for _, node := range graph.Nodes {
		if node.Name == name {
			return true
		}
	}

	return false
}

func (graph *Graph) AddAssociatedNodesTo(name string, associatedNodes []string) *Graph {
	for _, associatedNode := range associatedNodes {
		graph.AddAssociatedNodeTo(name, associatedNode)
	}
	return graph
}

func (graph *Graph) AddAssociatedNodeTo(name string, associatedNodeName string) *Graph {
	for i, node := range graph.Nodes {
		if node.Name == name {
			graph.Nodes[i].AssociatedNodesCount++
			graph.Nodes[i].AssociatedNodes = append(graph.Nodes[i].AssociatedNodes, associatedNodeName)
			return graph
		}
	}

	return graph.AddNode(Node{
		Name:                 name,
		AssociatedNodesCount: 1,
		AssociatedNodes:      []string{associatedNodeName},
	})
}

func (graph *Graph) GetNodeCount() int {
	return len(graph.Nodes)
}
