package analyzer

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const vertexPowerFile string = "VertexPower.log"

func (context *AnalyzerContext) GetVertexPower() map[string]int {
	return context.vertexPower
}

func (context *AnalyzerContext) SetVertextPower(vertexPower map[string]int) {
	context.vertexPower = vertexPower
}

func (context *AnalyzerContext) SaveContext(path string) {
	saveVertexPower(context, path)
}

func saveVertexPower(context *AnalyzerContext, path string) {
	file, err := os.Create(filepath.Join(path, vertexPowerFile))
	if err != nil {
		log.Fatal("Couldn't create ", vertexPowerFile, "by path ", path, err)
	}
	defer file.Close()

	log.Println("Saving power vertex results...")
	for vertex, power := range context.vertexPower {
		_, err = file.WriteString(vertex + ": " + strconv.Itoa(power) + "\n")
		if err != nil {
			log.Fatal("Error in the time of writing result to file")
		}
	}
}

type AnalyzerContext struct {
	vertexPower map[string]int
}
