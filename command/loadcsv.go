package command

import (
	"github.com/kshaposhnikov/twitter-crawler/crawler/user"
	"github.com/kshaposhnikov/twitter-crawler/graph"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2"
)

var loadCsvCmd = &cobra.Command{
	Use:   "load-csv",
	Short: "Load CSV file with data about Twitter users",
	Run:   loadCsv,
}

func loadCsv(cmd *cobra.Command, args []string) {
	if DB == nil {
		openDBConnection()
	}

	if len(args) < 1 {
		logrus.Fatal("Need to provide path to CSV file")
	}

	logrus.Debug("[Load CSV Command] Input Parameters", args)

	fillUsers(DB, args[0])
	users := user.FindAll(DB)
	buildGraph(DB, users)
}

func fillUsers(db *mgo.Database, path string) []user.User {
	users := user.ReadCSV(path)
	user.StoreUser(db, users...)
	return users
}

func buildGraph(db *mgo.Database, users []user.User) {
	for _, user := range users {
		graph.StoreVertex(db, graph.Node{
			Id:                   user.ID,
			AssociatedNodesCount: user.FriendsCount,
			AssociatedNodes:      user.Friends,
		})
	}
}
