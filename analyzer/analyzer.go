package analyzer

import (
	"github.com/kshaposhnikov/twitter-crawler/analyzer/graph"

	"gopkg.in/mgo.v2"
)

func CalcluatePowerByIter(iterator *mgo.Iter, context *AnalyzerContext) {
	powerVertex := make(map[string]int)
	var vertex graph.Vertex
	for iterator.Next(&vertex) {
		powerVertex[vertex.Name] += vertex.AssociatedVertexCount
		for _, vertex := range vertex.AssociatedVertexes {
			powerVertex[vertex]++
		}
	}
	context.vertexPower = powerVertex
	iterator.Close()
}

func CalculateProwerByArray(graph graph.Graph) map[string]int{
	powerVertex := make(map[string]int)
	for _, item := range graph {
		powerVertex[item.Name] += item.AssociatedVertexCount
		for _, vertex := range item.AssociatedVertexes {
			powerVertex[vertex]++
		}
	}
	return powerVertex
}

func CalculateDiameter() {

}
