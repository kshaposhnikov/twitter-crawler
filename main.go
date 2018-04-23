package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

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
	threadNumber := argFlags.Int("thread-number", 1, "Number of thread to asynchronous generation of graph")
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

			start := time.Now()
			log.Println(start)
			result := generator.SecondPhaseMultithreadGenerator{
				GeneralGenerator: generator.GeneralGenerator{
					VCount: n,
					ECount: m,
				},
				NumberThreads: *threadNumber,
			}.Generate()
			log.Println(time.Now().Sub(start))

			log.Println(">>> Final", result)

			// generator.GeneralGenerator{
			// 	VCount: n,
			// 	ECount: m,
			// }.Generate()
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
		graph.StoreVertex(db, graph.Node{
			Name:                 user.ID,
			AssociatedNodesCount: value,
			AssociatedNodes:      user.Friends,
		})
	}
}
