package command

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2"
	"log"
)

const skipDB = false

var DB *mgo.Database = nil

var rootCmd = &cobra.Command{
	Use: "twitter-crawler",
}

func init() {
	initLogger()
	rootCmd.AddCommand(loadCsvCmd, generateCmd, analyzeCmd, crawlerCmd)
}

func initLogger() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
}

func openDBConnection() {
	if skipDB {
		logrus.Info("Opening of DB connection was skipped")
		return
	}

	session, err := mgo.Dial("localhost")
	if err != nil {
		logrus.Error("[Root Command] DB Connection", err)
	} else {
		DB = session.DB("twitter-crawler")
	}
}

func closeConnection() {
	DB.Session.Close()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
