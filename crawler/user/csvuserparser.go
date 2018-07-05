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

	var arrays []arrayInterval
	var tmp arrayInterval
	for i := 0; i < len(columns); i++ {
		if strings.HasPrefix(columns[i], "[") {
			columns[i] = strings.TrimPrefix(columns[i], "[")
			tmp.start = i
		}

		if strings.HasSuffix(columns[i], "]") {
			columns[i] = strings.TrimSuffix(columns[i], "]")
			tmp.finish = i
			arrays = append(arrays, tmp)
		}

		columns[i] = strings.Trim(strings.TrimSpace(columns[i]), "\"")
	}

	return User{
		ID:             atoi(columns[0], 0),
		ScreenNames:    columns[1],
		Tags:           columns[arrays[0].start : arrays[0].finish],
		Avatar:         columns[arrays[0].finish + 1],
		FollowersCount: atoi(columns[arrays[0].finish + 2], 4),
		FriendsCount:   atoi(columns[arrays[0].finish + 3], 5),
		Lang:           columns[arrays[0].finish + 4],
		LastSeen:       columns[arrays[0].finish + 5],
		TweetID:        atoi(columns[arrays[0].finish + 6], 8),
		Friends: func() []int {
			var result []int
			for _, item := range columns[arrays[1].start : ] {
				result = append(result, atoi(item, 9))
			}
			return result
		}(),
	}
}

type arrayInterval struct {
	start int
	finish int
}

func atoi(str string, index int) int {
	if str == "" {
		return 0
	}

	res, err := strconv.Atoi(str)
	if err != nil {
		logrus.Fatalln("[csvuserparser.atoi] Cann't convert Str", str, "to int; Index: ", index)
	}
	return res
}
