package graph

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Gateway struct {
	db *mgo.Database
}

func NewGateway(db *mgo.Database) *Gateway {
	return &Gateway{
		db: db,
	}
}

// LoadGraph returns array of Graph structure by provided Mongo Database
func (self *Gateway) LoadGraph() Graph {
	log.Println("Loading graph...")
	var graph Graph
	err := self.db.C("graph").Find(bson.M{}).All(&graph)
	if err != nil {
		log.Fatal("Error in the time of loading all graph", err)
	}
	log.Println("Done.")
	return graph
}

func (self *Gateway) LoadGraphByIter() *mgo.Iter {
	return self.db.C("graph").Find(bson.M{}).Iter()
}

func (self *Gateway) Exists(userId int64) bool {
	//self.db.C("graph").Find(bson.)
	return true
}

func (self *Gateway) StoreVertex(node ...Node) {
	for iter, edge := range node {
		err := self.db.C("graph").Insert(edge)
		if err != nil {
			log.Fatal("#", iter, "Error in the time of insrting edge for node", edge.Id, "\n", err)
		}
	}
}
