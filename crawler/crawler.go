package crawler

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/sirupsen/logrus"
	"net/url"
	"strconv"
	"github.com/kshaposhnikov/twitter-crawler/graph"
	"gopkg.in/mgo.v2"
)

type Crawler struct {
	gateway    *graph.Gateway
	twitterApi *anaconda.TwitterApi
	depth      int
}

func New(api *anaconda.TwitterApi, db *mgo.Database, depthToLook int) *Crawler {
	return &Crawler{
		twitterApi: api,
		depth:      depthToLook,
		gateway:    graph.NewGateway(db),
	}
}

func (crw *Crawler) Start(startUser string) {
	users, err := crw.twitterApi.GetUserSearch(startUser, nil)
	if err != nil {
		logrus.Fatalln("Failed to search user", err)
	}

	crw.loadFollowers(users[0].Id, 0)
}

func (crw *Crawler) loadFollowers(userId int64, currentDepth int) {
	if currentDepth > crw.depth {
		return
	}

	v := url.Values{}
	v.Add("user_id", strconv.FormatInt(userId, 10))
	cursor, err := crw.twitterApi.GetFollowersIds(v)

	if err != nil {
		logrus.Error("Can't load followers for user_id", userId)
	}

	if !crw.gateway.Exists(userId) {
		crw.addToDataBase(userId, &cursor.Ids)
		return
	}

	currentDepth++
	for _, id := range cursor.Ids {
		crw.loadFollowers(id, currentDepth)
	}
}

func (crw *Crawler) addToDataBase(userId int64, followers *[]int64) {
	crw.gateway.StoreVertex(graph.Node{Id: userId})
	logrus.Info("User Id: ", userId, "\nFollowers Count: ", len(*followers), "\nFollowers: ", followers)
}
