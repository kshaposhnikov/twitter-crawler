package command

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/kshaposhnikov/twitter-crawler/crawler"
	"github.com/sirupsen/logrus"
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

func init() {
	crawlerCmd.Flags().StringVarP(&accessToken, "access-token", "", "", "Access Token from your account on Twitter Deveoper portal")
	crawlerCmd.Flags().StringVarP(&accessTokenSecret, "access-token-secret", "", "", "Access Token Secret from your account on Twitter Deveoper portal")
	crawlerCmd.Flags().StringVarP(&consumerKey, "consumer-key", "", "", "Consumer Key from your account on Twitter Deveoper portal")
	crawlerCmd.Flags().StringVarP(&consumerSecret, "consumer_secret", "", "", "Consumer Secret from your account on Twitter Deveoper portal")
	crawlerCmd.Flags().IntVarP(&depth, "depth", "d", 10, "Depth parameter to load")
}

func startCrawler(cmd *cobra.Command, args []string) {
	twitterApi := anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerSecret)
	defer twitterApi.Close()
	defer logrus.Info("Stop")

	//twitterApi.SetDelay(5 * time.Second)

	crawler.New(twitterApi, depth).Start("ShaposhnikovKs")

	time.Sleep(10 * time.Second)
}
