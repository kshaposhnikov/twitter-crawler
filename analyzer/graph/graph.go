package graph

type Vertex struct {
	Name                  string
	AssociatedVertexCount int
	AssociatedVertexes    []string
}

type Graph []Vertex

func (graph Graph) Concat(newGraph Graph) Graph{
	for _, vertex := range newGraph {
		graph = graph.AddVertex(vertex)
	}
	return graph
}

func (graph Graph) AddVertex(vertex Vertex) Graph {
	return append(graph, vertex)
}

func (graph Graph) ContainsVertex(name string) bool {
	for _, vertex := range graph {
		if vertex.Name == name {
			return true
		}
	}

	return false
}

func (graph Graph) AddAssociatedVertexesTo(name string, associatedVertexes []string) Graph {
	for _, associatedVertex := range associatedVertexes {
		graph = graph.AddAssociatedVertexTo(name, associatedVertex)
	}
	return graph
}

func (graph Graph) AddAssociatedVertexTo(name string, associatedVertexName string) Graph {
	for i, vertex := range graph {
		if vertex.Name == name {
			graph[i].AssociatedVertexCount++
			graph[i].AssociatedVertexes = append(graph[i].AssociatedVertexes, associatedVertexName)
			return graph
		}
	}

	return graph.AddVertex(Vertex{
		Name: name,
		AssociatedVertexCount: 1,
		AssociatedVertexes:    []string{associatedVertexName},
	})
}

func (graph Graph) GetVertexNumber() int {
	return len(graph)
}
