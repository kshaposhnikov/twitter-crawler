package generator

import "github.com/kshaposhnikov/twitter-crawler/graph"

type GeneralGenerator struct {
	VCount int
	ECount int
}

type Generator interface {
	Generate() graph.Graph
}

type FirstPhaseMultithreadGenerator struct {
	GeneralGenerator
	NumberThreads int
}

type SecondPhaseMultithreadGenerator struct {
	GeneralGenerator
	NumberThreads int
}
