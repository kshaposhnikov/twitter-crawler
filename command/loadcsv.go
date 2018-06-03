package command

import (
	"github.com/kshaposhnikov/twitter-crawler/crawler/user"
	"github.com/kshaposhnikov/twitter-crawler/graph"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

	fillUsers(args[0])
	users := user.FindAll(DB)
	buildGraph(users)
}

func fillUsers(path string) []user.User {
	users := user.ReadCSV(path)
	user.StoreUser(DB, users...)
	return users
}

func buildGraph(users []user.User) {
	gateway := graph.NewGateway(DB)
	for _, user := range users {
		gateway.StoreVertex(graph.Node{
			Id:                   user.ID,
			AssociatedNodesCount: user.FriendsCount,
			AssociatedNodes:      user.Friends,
		})
	}
}
