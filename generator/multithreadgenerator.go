package generator

import (
	"github.com/kshaposhnikov/twitter-crawler/graph"
	"log"
)

func (gen FirstPhaseMultithreadGenerator) Generate() graph.Graph {
	n := (gen.VCount * gen.ECount) / gen.NumberThreads
	graphs := make(chan *graph.Graph)
	for i := 0; i < gen.NumberThreads; i++ {
		go buildInitialGraph(n, graphs)
	}

	intermidiateGraph := graph.NewGraph()
	for i := 0; i < gen.NumberThreads; i++ {
		intermidiateGraph = intermidiateGraph.Concat(<-graphs)
	}

	log.Println(*intermidiateGraph)

	j := n
	for i := n; i < intermidiateGraph.GetNodeCount(); i++ {
		intermidiateGraph.Nodes[i].Id += int64(j)

		for l := 0; l < len(intermidiateGraph.Nodes[i].AssociatedNodes); l++ {
			intermidiateGraph.Nodes[i].AssociatedNodes[l] += int64(j)
		}

		if i == j*2-1 {
			j *= 2
		}
	}

	log.Println(intermidiateGraph)

	return buildFinalGraph(intermidiateGraph, gen.ECount)
}

func buildInitialGraph(n int, graphs chan *graph.Graph) {
	simpleGenerator := GeneralGenerator{}
	graphs <- simpleGenerator.buildInitialGraph(n)
}

func buildFinalGraph(initialGraph *graph.Graph, m int) graph.Graph {
	simpleGenerator := GeneralGenerator{}
	return simpleGenerator.buildFinalGraph(initialGraph, 0, initialGraph.GetNodeCount(), int64(m))
}
