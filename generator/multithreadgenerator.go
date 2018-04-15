package generator

import (
	"github.com/kshaposhnikov/twitter-crawler/analyzer/graph"
	"log"
	"strconv"
)

type MultithreadGenerator struct {
	GeneralGenerator
	NumberThreads int
}

func (gen MultithreadGenerator) Generate() graph.Graph {
	n := (gen.VCount * gen.ECount) / gen.NumberThreads
	graphs := make(chan graph.Graph)
	for i := 0; i < gen.NumberThreads; i++ {
		go buildInitialGraph(n, graphs)
	}

	intermidiateGraph := graph.NewGraph()
	for i := 0; i < gen.NumberThreads; i++ {
		intermidiateGraph = intermidiateGraph.Concat(<- graphs)
	}

	log.Println(*intermidiateGraph)

	j := n
	for i := n; i < intermidiateGraph.GetNodeCount(); i++ {
		intermidiateGraph.Nodes[i].Name = updateVertex(intermidiateGraph.Nodes[i].Name, j)

		for l := 0; l < len(intermidiateGraph.Nodes[i].AssociatedNodes); l++ {
			intermidiateGraph.Nodes[i].AssociatedNodes[l] = updateVertex(intermidiateGraph.Nodes[i].AssociatedNodes[l], j)
		}

		if i == j * 2 - 1 {
			j *= 2
		}
	}

	log.Println(intermidiateGraph)

	result := buildFinalGraph(intermidiateGraph, gen.ECount)

	return result
}

func updateVertex(vertex string, j int) string {
	current, _ := strconv.Atoi(vertex)
	return strconv.Itoa(current + j)
}

func buildInitialGraph(n int, graphs chan graph.Graph) {
	simpleGenerator := GeneralGenerator{}
	graphs <- simpleGenerator.buildInitialGraph(n)
}

func buildFinalGraph(initialGraph *graph.Graph, m int) graph.Graph {
	simpleGenerator := GeneralGenerator{}
	return simpleGenerator.buildFinalGraph(initialGraph, m)
}