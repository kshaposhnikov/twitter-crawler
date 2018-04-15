package analyzer

import (
	"github.com/kshaposhnikov/twitter-crawler/analyzer/graph"

	"gopkg.in/mgo.v2"
)

func CalcluatePowerByIter(iterator *mgo.Iter, context *AnalyzerContext) {
	powerVertex := make(map[string]int)
	var vertex graph.Node
	for iterator.Next(&vertex) {
		powerVertex[vertex.Name] += vertex.AssociatedNodesCount
		for _, vertex := range vertex.AssociatedNodes {
			powerVertex[vertex]++
		}
	}
	context.vertexPower = powerVertex
	iterator.Close()
}

func CalculateProwerByArray(graph *graph.Graph) map[string]int{
	powerVertex := make(map[string]int)
	for _, item := range graph.Nodes {
		powerVertex[item.Name] += item.AssociatedNodesCount
		for _, vertex := range item.AssociatedNodes {
			powerVertex[vertex]++
		}
	}
	return powerVertex
}

func CalculateDiameter() {

}
