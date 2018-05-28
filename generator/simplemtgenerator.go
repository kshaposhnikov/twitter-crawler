package generator

import (
	"github.com/kshaposhnikov/twitter-crawler/graph"
	"sync"
)

func (gen SecondPhaseMultithreadGenerator) Generate() graph.Graph {
	generator := GeneralGenerator{
		gen.VCount,
		gen.ECount,
	}
	initialGraph := generator.buildInitialGraph(gen.VCount * gen.ECount)
	batch := calculateInterval(gen.VCount*gen.ECount, gen.NumberThreads)
	goroutineNumber := calculateInterval(initialGraph.GetNodeCount(), batch)
	graphs := make(chan graph.Graph, goroutineNumber)
	var wg sync.WaitGroup
	wg.Add(goroutineNumber)
	for i := 0; i < goroutineNumber; i++ {
		from := i * batch
		to := from + batch
		if to >= initialGraph.GetNodeCount() {
			to = initialGraph.GetNodeCount()
		}
		go func() {
			defer wg.Done()
			graphs <- generator.buildFinalGraph(initialGraph, from, to, gen.ECount)
		}()
	}
	wg.Wait()
	close(graphs)

	result := graph.NewGraph()
	for item := range graphs {
		result.Concat(&item)
	}

	return *result
}
