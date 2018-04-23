package generator

import "github.com/kshaposhnikov/twitter-crawler/analyzer/graph"
import (
	//"github.com/kshaposhnikov/twitter-crawler/analyzer"
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"

	"gonum.org/v1/gonum/floats"
)

//bollobas-riordan
// Number of threads should be less then m
func (gen GeneralGenerator) Generate() graph.Graph {
	n := gen.VCount
	m := gen.ECount
	if m < 2 {
		log.Fatal("m should more or equal 2")
	}

	var previousGraph = gen.buildInitialGraph(n * m)
	return gen.buildFinalGraph(&previousGraph, 0, previousGraph.GetNodeCount(), 1, m)
}

func (gen GeneralGenerator) buildInitialGraph(n int) graph.Graph {
	previousGraph := graph.NewGraph()
	previousGraph.AddNode(graph.Node{
		Name:                 "1",
		AssociatedNodesCount: 1,
		AssociatedNodes:      []string{"1"},
	})

	degree := make(map[int]int)
	degree[0] = 2

	for i := 1; i <= n-1; i++ {
		previousGraph = nextGraph(previousGraph, degree)
		//log.Println(">>> Initial Graph", previousGraph)
	}

	return *previousGraph
}

func nextGraph(previousGraph *graph.Graph, degrees map[int]int) *graph.Graph {
	//random := rand.New(rand.NewSource(time.Now().UnixNano()))

	probabilities := calculateProbabilities(degrees)
	cdf := cumsum(probabilities)

	x := rand.Float64()
	idx := sort.Search(len(cdf), func(i int) bool {
		return cdf[i] > x
	})

	var node graph.Node
	if idx > previousGraph.GetNodeCount() {
		node = graph.Node{
			Name:                 strconv.Itoa(len(probabilities)),
			AssociatedNodesCount: 1,
			AssociatedNodes:      []string{strconv.Itoa(len(probabilities))},
		}
		degrees[len(probabilities)-1]++
	} else {
		degrees[idx]++
		node = graph.Node{
			Name:                 strconv.Itoa(len(probabilities)),
			AssociatedNodesCount: 1,
			AssociatedNodes:      []string{strconv.Itoa(idx + 1)},
		}
	}

	degrees[len(probabilities)-1]++
	return previousGraph.AddNode(node)
}

func (gen GeneralGenerator) buildFinalGraph(pregeneratedGraph *graph.Graph, from, to, left, m int) graph.Graph {
	result := graph.NewGraph()

	j := left/m + 1
	right := left + m
	loops := []string{}
	for i, node := range pregeneratedGraph.Nodes[from:to] {
		for _, associatedVertex := range node.AssociatedNodes {
			current, _ := strconv.Atoi(associatedVertex)
			if current < right && current > left {
				loops = append(loops, strconv.Itoa(j))
			} else if current >= right || current <= left {
				result = result.AddAssociatedNodeTo(strconv.Itoa(j), strconv.Itoa(calculateInterval(current, m)))
			}
		}

		if i+from == right-1 {
			if len(loops) > 0 {
				result = result.AddAssociatedNodesTo(strconv.Itoa(j), loops)
			} else if !result.ContainsVertex(strconv.Itoa(j)) {
				result = result.AddNode(graph.Node{
					Name:                 strconv.Itoa(j),
					AssociatedNodesCount: len(loops),
					AssociatedNodes:      loops,
				})
			}
			loops = []string{}
			left = right
			right += m
			j++
		}
	}

	return *result
}

func calculateInterval(number int, m int) int {
	if number%m == 0 {
		return number / m
	} else {
		return int(math.Trunc(float64(number)/float64(m)) + 1)
	}
}

func calculateProbabilities(degrees map[int]int) []float64 {
	n := float64(len(degrees) + 1)
	var probabilities []float64
	for _, power := range degrees {
		probabilities = append(probabilities, float64(power)/(2.0*n-1.0))
	}
	return append(probabilities, 1.0/(2.0*n-1.0))
}

func cumsum(probabilities []float64) []float64 {
	dest := make([]float64, len(probabilities))
	return floats.CumSum(dest, probabilities)
}
