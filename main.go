// dailygo project main.go
package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/tucnak/telebot"
)

type Config struct {
	Token    string
	Interval time.Duration
	AdminID  int
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

	bot, err := telebot.NewBot(conf.Token)
	if err != nil {
		log.Fatal(err)
	}

	messages := make(chan telebot.Message)
	bot.Listen(messages, conf.Interval*time.Second)

	for message := range messages {
		if message.Text == "/start" {
			bot.SendMessage(message.Chat,
				"Hello, "+message.Sender.FirstName+"!", nil)
		}
	}
}
