package generator

import "github.com/kshaposhnikov/twitter-crawler/analyzer/graph"
import (
	//"github.com/kshaposhnikov/twitter-crawler/analyzer"
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"

	"runtime"
	"sync"
	"time"

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

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	degree := make(map[int]int)
	degree[0] = 2
	for i := 1; i <= n-1; i++ {
		previousGraph = nextGraph(previousGraph, degree, random)
		//log.Println(">>> Initial Graph", previousGraph)
	}

	return *previousGraph
}

func nextGraph(previousGraph *graph.Graph, degrees map[int]int, random *rand.Rand) *graph.Graph {
	probabilities := mtCalculateProbabilities(degrees)
	//probabilities := calculateProbabilities(degrees, 0, len(degrees))
	cdf := cumsum(probabilities)

	x := random.Float64()
	idx := sort.Search(len(cdf), func(i int) bool {
		return cdf[i] > x
	})

	// var node graph.Node
	// if idx > previousGraph.GetNodeCount() {
	// 	node = graph.Node{
	// 		Name:                 strconv.Itoa(len(probabilities)),
	// 		AssociatedNodesCount: 1,
	// 		AssociatedNodes:      []string{strconv.Itoa(len(probabilities))},
	// 	}
	// 	degrees[len(probabilities)-1]++
	// } else {
	// 	degrees[idx]++
	// 	node = graph.Node{
	// 		Name:                 strconv.Itoa(len(probabilities)),
	// 		AssociatedNodesCount: 1,
	// 		AssociatedNodes:      []string{strconv.Itoa(idx + 1)},
	// 	}
	// }

	degrees[idx]++

	degrees[len(probabilities)-1]++
	return previousGraph.AddNode(graph.Node{
		Name:                 strconv.Itoa(len(probabilities)),
		AssociatedNodesCount: 1,
		AssociatedNodes:      []string{strconv.Itoa(idx + 1)},
	})
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

func mtCalculateProbabilities(degrees map[int]int) []float64 {
	cpu := runtime.NumCPU()
	if len(degrees) > cpu*10 {
		probabilityResults := make(chan probabilityResult, 5)
		batch := int(math.Trunc(float64(len(degrees))/float64(cpu))) + 1
		var wg sync.WaitGroup
		wg.Add(cpu)
		for i := 0; i < cpu; i++ {
			from := i * batch
			to := from + batch
			if i == cpu-1 {
				to = len(degrees) - 2
			}

			go func(degrees map[int]int, from, to, order int, res chan probabilityResult) {
				defer wg.Done()
				res <- probabilityResult{
					order,
					calculateProbabilities(degrees, from, to),
				}
			}(degrees, from, to, i, probabilityResults)
		}
		wg.Wait()
		close(probabilityResults)

		var probabilities []probabilityResult
		for result := range probabilityResults {
			probabilities = append(probabilities, result)
		}
		sort.Slice(probabilities, func(i, j int) bool {
			return probabilities[i].order > probabilities[j].order
		})
		var result []float64
		for _, item := range probabilities {
			result = append(result, item.probabilities...)
		}
		return result
	} else {
		return calculateProbabilities(degrees, 0, len(degrees))
	}
}

func calculateProbabilities(degrees map[int]int, from, to int) []float64 {
	n := float64(len(degrees) + 1)
	var probabilities []float64
	for i := from; i < to; i++ {
		probabilities = append(probabilities, float64(degrees[i])/(2.0*n-1.0))
	}
	return append(probabilities, 1.0/(2.0*n-1.0))
}

func cumsum(probabilities []float64) []float64 {
	dest := make([]float64, len(probabilities))
	return floats.CumSum(dest, probabilities)
}

type probabilityResult struct {
	order         int
	probabilities []float64
}
