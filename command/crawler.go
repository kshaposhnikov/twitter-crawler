package command

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/kshaposhnikov/twitter-crawler/crawler"
	//"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var crawlerCmd = &cobra.Command{
	Use:   "start-crawler",
	Short: "Start crawler to build graph by Twitter",
	Run:   startCrawler,
}

var accessToken string
var accessTokenSecret string
var consumerKey string
var consumerSecret string
var depth int
var startFromUser string
var startFromUserId int64

func init() {
	crawlerCmd.Flags().StringVarP(&accessToken, "access-token", "", "", "Access Token from your account on Twitter Deveoper portal")
	crawlerCmd.Flags().StringVarP(&accessTokenSecret, "access-token-secret", "", "", "Access Token Secret from your account on Twitter Deveoper portal")
	crawlerCmd.Flags().StringVarP(&consumerKey, "consumer-key", "", "", "Consumer Key from your account on Twitter Deveoper portal")
	crawlerCmd.Flags().StringVarP(&consumerSecret, "consumer_secret", "", "", "Consumer Secret from your account on Twitter Deveoper portal")
	crawlerCmd.Flags().StringVarP(&startFromUser, "start_from_user", "", "", "User name that will be used as starting point")
	crawlerCmd.Flags().Int64VarP(&startFromUserId, "start_from_user_id", "", 0, "User Id that will be used as starting point")
	crawlerCmd.Flags().IntVarP(&depth, "depth", "d", 1, "Depth parameter to load")
}

func startCrawler(cmd *cobra.Command, args []string) {
	twitterApi := anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerSecret)
	defer twitterApi.Close()

	twitterApi.EnableThrottling(2*time.Second, 1)

	if DB == nil {
		openDBConnection()
		defer closeConnection()
	}

	crw := crawler.New(twitterApi, DB, depth)
	if startFromUser != "" {
		crw.StartByName(startFromUser)
	} else {
		crw.StartById(startFromUserId)
	}
}
