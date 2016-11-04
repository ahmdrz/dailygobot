// dailygo project main.go
package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ahmdrz/goinsta"
)

type Config struct {
	Username string
	Password string
}

var users map[string]bool = make(map[string]bool)

func main() {
	var tomlData string
	if data, err := ioutil.ReadFile("config.toml"); err != nil {
		log.Fatal(err)
	} else {
		tomlData = string(data)
	}

	var conf Config
	if _, err := toml.Decode(tomlData, &conf); err != nil {
		log.Fatal(err)
	}

	log.Println("Connecting to instagram server ...")
	insta := goinsta.New(conf.Username, conf.Password)
	err := insta.Login()
	if err != nil {
		log.Fatal(err)
	}
	feeds, err := insta.TagFeed("Programmer")
	if err != nil {
		log.Fatal(err)
	}
	for _, feed := range feeds.Items {
		if _, ok := users[feed.User.StringID()]; ok {
			continue
		}
		users[feed.User.StringID()] = true
		if feed.HasLiked {
			continue
		}

		if feed.LikeCount > 50 {
			_, err := insta.Comment(feed.ID, "Daily Golang contest on my page. Follow me , and if you want to get more contents join our Telegram channel. (sent by instagram bot) ")
			if err != nil {
				log.Println(feed.ID, "Comment failed")
			} else {
				log.Println(feed.ID, "Comment done")
			}
		}

		_, err := insta.Like(feed.ID)
		if err != nil {
			log.Println(feed.ID, "Failed")
		} else {
			log.Println(feed.ID, "Liked")
		}
		time.Sleep(1 * time.Second)
	}
}
