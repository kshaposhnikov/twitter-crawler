package crawler

import (
	"context"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
)

func GetClient(accessToken *string) *twitter.Client {
	config := oauth2.Config{}
	token := &oauth2.Token{AccessToken: *accessToken}
	client := config.Client(context.Background(), token)
	return twitter.NewClient(client)
}

func GetTwitterUser(twitterClient *twitter.Client, screenName string) {
	userParams := twitter.UserLookupParams{ScreenName: []string{screenName}}
	users, _, _ := twitterClient.Users.Lookup(&userParams)

	for _, twitterUser := range users {
		fmt.Println("Current user has ", twitterUser.FollowersCount)
	}
}
