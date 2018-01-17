package graph

import (
	"log"
	"os"
)

type GraphDBConfig struct {
	CollectionName string
}

func (config *GraphDBConfig) fillConfig(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Couldn't open provided file ", path, err)
	}
	defer file.Close()

}
