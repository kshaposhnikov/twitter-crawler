package user

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID             string
	ScreenNames    string
	Tags           string
	Avatar         string
	FollowersCount string
	FriendsCount   string
	Lang           string
	LastSeen       string
	TweetID        string
	Friends        []string
}

func FindAll(database *mgo.Database) []User {
	log.Println("Loading users...")
	var users []User
	err := database.C("users").Find(bson.M{}).All(&users)
	if err != nil {
		log.Fatal("Error in the time of finding all users", err)
	}
	log.Println("Done.")
	return users
}

func StoreUser(database *mgo.Database, users ...User) {
	for iter, user := range users {
		err := database.C("users").Insert(user)
		if err != nil {
			log.Fatal("#", iter, "Error in the time of insrting user ", user.ID, "\n", err)
		}
	}
}
