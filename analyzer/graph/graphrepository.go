package graph

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// LoadGraph returns array of Graph structure by provided Mongo Database
func LoadGraph(database *mgo.Database) Graph {
	log.Println("Loading graph...")
	var graph Graph
	err := database.C("graph").Find(bson.M{}).All(&graph)
	if err != nil {
		log.Fatal("Error in the time of loading all graph", err)
	}
	log.Println("Done.")
	return graph
}

func LoadGraphByIter(database *mgo.Database) *mgo.Iter {
	return database.C("graph").Find(bson.M{}).Iter()
}

func StoreVertex(database *mgo.Database, node ...Node) {
	for iter, edge := range node {
		err := database.C("graph").Insert(edge)
		if err != nil {
			log.Fatal("#", iter, "Error in the time of insrting edge for node", edge.Name, "\n", err)
		}
	}
}
