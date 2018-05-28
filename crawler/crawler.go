package crawler

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/sirupsen/logrus"
	"net/url"
)

type Crawler struct {
	twitterApi *anaconda.TwitterApi
	depth      int
}

func New(api *anaconda.TwitterApi, depthToLook int) *Crawler {
	return &Crawler{
		twitterApi: api,
		depth:      depthToLook,
	}
}

func (crw *Crawler) Start(startUser string) {
	users, err := crw.twitterApi.GetUserSearch(startUser, nil)
	if err != nil {
		logrus.Fatalln("Failed to search user", err)
	}

	loadFollowers(users[0].Id, 0)
}

func (crw *Crawler) loadFollowers(userId int, currentDepth int) {
	if currentDepth == crw.depth {
		return
	}

	v := url.Values{"user_id": userId}
	for page := range crw.twitterApi.GetFollowersIdsAll(v) {
		currentDepth++

	}

	for _, user := range users {
		crw.twitterApi.GetFollowersIdsAll(v)
		logrus.Info(user.Name, " ", user.FollowersCount)
	}
}
