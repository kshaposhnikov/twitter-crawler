package analyzer

import (
	"github.com/kshaposhnikov/twitter-crawler/graph"

	"gopkg.in/mgo.v2"
)

func CalcluatePowerByIter(iterator *mgo.Iter, context *AnalyzerContext) {
	powerVertex := make(map[int64]int)
	var vertexTmp graph.Node
	for iterator.Next(&vertexTmp) {
		powerVertex[vertexTmp.Id] += vertexTmp.AssociatedNodesCount
		for _, vertex := range vertexTmp.AssociatedNodes {
			powerVertex[vertex]++
		}
	}
	context.vertexPower = powerVertex
	iterator.Close()
}

func CalculateProwerByArray(graph *graph.Graph) map[int64]int {
	powerVertex := make(map[int64]int)
	for _, item := range graph.Nodes {
		powerVertex[item.Id] += item.AssociatedNodesCount
		for _, vertex := range item.AssociatedNodes {
			powerVertex[vertex]++
		}
	}
	return powerVertex
}

func CalculateDiameter() {

}
