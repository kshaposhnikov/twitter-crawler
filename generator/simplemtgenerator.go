package generator

import "github.com/kshaposhnikov/twitter-crawler/analyzer/graph"
import "log"

func (gen SecondPhaseMultithreadGenerator) Generate() graph.Graph {
	generator := GeneralGenerator{
		gen.VCount,
		gen.ECount,
	}
	initialGraph := generator.buildInitialGraph(gen.VCount * gen.ECount)
	graphs := make(chan graph.Graph)
	batch := (gen.VCount * gen.ECount) / gen.NumberThreads
	for i := 0; i < gen.NumberThreads; i++ {
		left := i * batch
		go execInNewThread(&initialGraph, left, left+batch, left, gen.ECount, graphs)
	}

	for i := 0; i < gen.NumberThreads; i++ {
		log.Println("Final", <-graphs)
	}

	return *graph.NewGraph()
}

func execInNewThread(initialGraph *graph.Graph, from, to, left, m int, graphs chan graph.Graph) {
	simpleGenerator := GeneralGenerator{}
	graphs <- simpleGenerator.buildFinalGraph(initialGraph, from, to, left, m)
}
