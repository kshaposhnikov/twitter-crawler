package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"gopkg.in/mgo.v2"

	"regexp"
	"strings"

	"github.com/kshaposhnikov/twitter-crawler/analyzer"
	"github.com/kshaposhnikov/twitter-crawler/analyzer/graph"
	"github.com/kshaposhnikov/twitter-crawler/crawler/user"
	"github.com/kshaposhnikov/twitter-crawler/generator"
)

func main() {
	argFlags := flag.NewFlagSet("twitter-crawler", flag.ExitOnError)
	//accessToken := argFlags.String("twitter-access-token", "", "Twitter Access Token")
	csvFile := argFlags.String("load-csv", "", "CSV file with data about Twitter users")
	analyze := argFlags.Bool("analyze", false, "Start analyzing of existing graph")
	generate := argFlags.String("generate", "", "Generate new graph")
	argFlags.Parse(os.Args[1:])

	var db *mgo.Database = nil

	if *csvFile != "" {
		db = openConnection()
		fillUsers(db, *csvFile)
		users := user.FindAll(db)
		buildGraph(db, users)
	}

	if *analyze {
		if db == nil {
			db = openConnection()
		}
		graphIterator := graph.LoadGraphByIter(db)
		var context analyzer.AnalyzerContext
		analyzer.CalcluatePowerByIter(graphIterator, &context)
		context.SaveContext("C:\\analyze graph\\")
	}

	if *generate != "" {
		graphConfig := regexp.MustCompile(`[0-9]+;[0-9]+`)
		if graphConfig.MatchString(*generate) {
			nm := strings.Split(graphConfig.FindString(*generate), ";")
			n, _ := strconv.Atoi(nm[0])
			m, _ := strconv.Atoi(nm[1])

			generator.MultithreadGenerator{
					GeneralGenerator: generator.GeneralGenerator{
					VCount: n,
					ECount: m,
				},
				NumberThreads: 3,
			}.Generate()
		} else {
			log.Fatalln("Need to specify format `n;m`")
		}
	}
}

func openConnection() *mgo.Database {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	} else {
		defer session.Close()
		return session.DB("twitter-crawler")
	}

	return nil
}

func fillUsers(db *mgo.Database, path string) []user.User {
	users := user.ReadCSV(path)
	user.StoreUser(db, users...)
	return users
}

func buildGraph(db *mgo.Database, users []user.User) {
	for _, user := range users {
		value, _ := strconv.Atoi(user.FriendsCount)
		graph.StoreVertex(db, graph.Vertex{
			Name: user.ID,
			AssociatedVertexCount: value,
			AssociatedVertexes:    user.Friends,
		})
	}
}
