package user

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func ReadCSV(path string) []User {
	data, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReaderSize(data, 512000)
	var parsedData []User
	for i := 0; ; i++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Fatal("Error in the time of reading", err)
			} else {
				break
			}
		}
		log.Print("Start parsing line ", i)
		parsedData = append(parsedData, parseCSVLine(string(line)))
	}
	return parsedData
}

func parseCSVLine(inputLine string) User {
	columns := strings.Split(inputLine, ",")
	log.Print("Parse info for user ", columns[0])
	if strings.HasPrefix(columns[9], "[") {
		columns[9] = strings.TrimPrefix(columns[9], "[")
		columns[len(columns)-1] = strings.TrimSuffix(columns[len(columns)-1], "]")
	}

	for i := 0; i < len(columns); i++ {
		columns[i] = strings.Trim(strings.TrimSpace(columns[i]), "\"")
	}

	return User{
		ID:             atoi(columns[0]),
		ScreenNames:    columns[1],
		Tags:           columns[2],
		Avatar:         columns[3],
		FollowersCount: atoi(columns[4]),
		FriendsCount:   atoi(columns[5]),
		Lang:           columns[6],
		LastSeen:       columns[7],
		TweetID:        atoi(columns[8]),
		Friends: func() []int {
			result := make([]int, len(columns[9:]))
			for _, item := range columns[9:] {
				result = append(result, atoi(item))
			}
			return result
		}(),
	}
}

func atoi(str string) int {
	res, err := strconv.Atoi(str)
	if err != nil {
		logrus.Fatalln("Cann't convert Str", str, "to int")
	}
	return res
}
