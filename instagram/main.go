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
	feeds, err := insta.TagFeed("golang")
	if err != nil {
		log.Fatal(err)
	}
	for _, feed := range feeds.Items {
		_, err := insta.Like(feed.ID)
		if err != nil {
			log.Println(feed.ID, "Failed")
		} else {
			log.Println(feed.ID, "Liked")
		}
		time.Sleep(1 * time.Second)
	}
}
