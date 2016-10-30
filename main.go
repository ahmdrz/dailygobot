// dailygo project main.go
package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tucnak/telebot"
)

type Config struct {
	Token    string
	Interval time.Duration
	Admin    int64
}

type User struct {
	UserId   int
	Position int
	LastUsed int64
}

func main() {
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
	}

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

	log.Println("Connecting to telegram server ...")
	bot, err := telebot.NewBot(conf.Token)
	if err != nil {
		log.Fatal(err)
	}

	messages := make(chan telebot.Message)
	bot.Listen(messages, conf.Interval*time.Second)

	log.Printf("Listening on @%s ...", bot.Identity.Username)

	for message := range messages {
		log.Println(message.Sender.ID, message.Sender.Username, message.Text)

		user := User{}
		db.Table("users").Find(&user, "user_id = ?", message.Sender.ID)
		if user.UserId != message.Sender.ID {
			log.Println("not found")
			user.UserId = message.Sender.ID
		}
		user.LastUsed = time.Now().Unix()

		if message.Text == "/start" {
			user.Position = 0
			bot.SendMessage(message.Sender, "Hello, "+message.Sender.FirstName+"! I'm the daily Go bot ,\nChoose an option.", &telebot.SendOptions{
				ReplyMarkup: telebot.ReplyMarkup{
					Selective:  true,
					ForceReply: true,
					CustomKeyboard: [][]string{
						[]string{"Suggest"},
						[]string{"Contact to Admin"},
					},
					ResizeKeyboard:  true,
					OneTimeKeyboard: true,
				},
				ReplyTo: message,
			})
		} else if message.Text == "Suggest" {
			user.Position = 1
			bot.SendMessage(message.Sender, "Enter your suggestion", nil)
		} else if message.Text == "Contact to Admin" {
			user.Position = 2
			bot.SendMessage(message.Sender, "Enter your message for admin", nil)
		} else {
			if user.Position == 1 {
				err = bot.SendMessage(telebot.Chat{ID: conf.Admin}, "Suggest : "+message.Text+" - @"+message.Sender.Username+" ["+message.Sender.FirstName+" "+message.Sender.LastName+"] ", nil)
				if err != nil {
					log.Println(err.Error())
					bot.SendMessage(message.Sender, "Can't send message", nil)
				} else {
					bot.SendMessage(message.Sender, "Thanks for your suggestion", &telebot.SendOptions{
						ReplyMarkup: telebot.ReplyMarkup{
							Selective:  true,
							ForceReply: true,
							CustomKeyboard: [][]string{
								[]string{"Suggest"},
								[]string{"Contact to Admin"},
							},
							ResizeKeyboard:  true,
							OneTimeKeyboard: true,
						},
						ReplyTo: message,
					})
				}
			} else if user.Position == 2 {
				err = bot.SendMessage(telebot.Chat{ID: conf.Admin}, "Contact : "+message.Text+" - @"+message.Sender.Username+" ["+message.Sender.FirstName+" "+message.Sender.LastName+"] ", nil)
				if err != nil {
					log.Println(err.Error())
					bot.SendMessage(message.Sender, "Can't send message", nil)
				} else {
					bot.SendMessage(message.Sender, "Message has been sent to admin", &telebot.SendOptions{
						ReplyMarkup: telebot.ReplyMarkup{
							Selective:  true,
							ForceReply: true,
							CustomKeyboard: [][]string{
								[]string{"Suggest"},
								[]string{"Contact to Admin"},
							},
							ResizeKeyboard:  true,
							OneTimeKeyboard: true,
						},
						ReplyTo: message,
					})
				}
			}
		}

		db.Table("users").Save(&user)
	}
}
