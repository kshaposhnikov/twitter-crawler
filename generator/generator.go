package generator

import "github.com/kshaposhnikov/twitter-crawler/analyzer/graph"

type GeneralGenerator struct {
	VCount int
	ECount int
}

type Generator interface {
	Generate() graph.Graph
}

type BollobasRiordanGenerator interface {
	Generator
	buildInitialGraph(n int) graph.Graph
	buildFinalGraph(initialGraph graph.Graph, m int) graph.Graph
}