package command

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

const skipDB = true

var DB *mgo.Database = nil

var rootCmd = &cobra.Command{
	Use: "twitter-crawler",
}

func init() {
	initLogger()
	rootCmd.AddCommand(loadCsvCmd, generateCmd, analyzeCmd)
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
		defer session.Close()
		DB = session.DB("twitter-crawler")
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
