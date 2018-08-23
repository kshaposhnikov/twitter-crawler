package crawler

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/kshaposhnikov/twitter-crawler/graph"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"net/url"
	"strconv"
	"time"
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

func (crw *Crawler) StartByName(startUser string) {
	users, err := crw.twitterApi.GetUserSearch(startUser, nil)
	if err != nil {
		logrus.Fatalln("Failed to search user", err)
	}

	crw.StartById(users[0].Id)
}

func (crw *Crawler) StartById(userId int64) {
	crw.loadFollowers(userId, 0)
}

func (crw *Crawler) loadFollowers(userId int64, currentDepth int) {
	if currentDepth > crw.depth {
		return
	}

	wait(1)

	v := url.Values{}
	v.Add("user_id", strconv.FormatInt(userId, 10))
	logrus.Debug("Get friends for ", userId)
	cursor, err := crw.twitterApi.GetFriendsIds(v)

	if err != nil {
		logrus.Error("Can't load followers for user_id", userId)
	}

	if !crw.gateway.Exists(userId) {
		crw.addToDataBase(userId, &cursor.Ids)
	} else {
		logrus.Info("User Id: ", userId, " already exists")
	}

	currentDepth++
	for _, id := range cursor.Ids {
		crw.loadFollowers(id, currentDepth)
	}
}

func (crw *Crawler) addToDataBase(userId int64, friends *[]int64) {
	crw.gateway.StoreVertex(graph.Node{
		Id: userId,
		AssociatedNodesCount: len(*friends),
		AssociatedNodes: *friends,

	})
	logrus.Info("User Id: ", userId, "\nFriends Count: ", len(*friends), "\nFriends: ", friends)
}

func wait(period time.Duration) {
	timer := time.NewTimer(period * time.Minute)
	<- timer.C
	timer.Stop()
}