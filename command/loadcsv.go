package command

import (
	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
	"github.com/kshaposhnikov/twitter-crawler/crawler/user"
	"gopkg.in/mgo.v2"
	"strconv"
	"github.com/kshaposhnikov/twitter-crawler/analyzer/graph"
)

var loadCsvCmd = &cobra.Command{
	Use: "load-csv",
	Short: "Load CSV file with data about Twitter users",
	Run: loadCsv,
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
		value, _ := strconv.Atoi(user.FriendsCount)
		graph.StoreVertex(db, graph.Node{
			Name:                 user.ID,
			AssociatedNodesCount: value,
			AssociatedNodes:      user.Friends,
		})
	}
}