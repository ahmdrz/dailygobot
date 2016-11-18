package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/tucnak/telebot"
)

var keyboard [][]string = [][]string{
	[]string{"New Link"},
	[]string{"New Post"},
	[]string{"Get Link Info"},
	[]string{"Post To Site"},
}

const (
	TOKEN   = "268916701:AAHzujvtmANyixZ_wY1Q8kJPRtfPO2Ru6Ts" // sample token !
	PORT    = ":8090"
	ADMINID = 83919508
)

func main() {
	open_db()
	rand.Seed(time.Now().Unix())

	bot, err := telebot.NewBot(TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	go telegram(bot)

	router := mux.NewRouter()
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("golng-static/static/"))))
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/{id}", ShortnerHandler)
	router.HandleFunc("/api/{time}", ApiHandler)
	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(PORT, handler))
}

func telegram(bot *telebot.Bot) {
	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {
		msg := message.Text
		user := User{}
		db.Table("users").Find(&user, "user_id = ?", message.Sender.ID)
		if user.UserId != message.Sender.ID {
			user.UserId = message.Sender.ID
		}
		user.LastUsed = time.Now().Unix()

		if message.Sender.ID == ADMINID {
			if msg == "/start" {
				user.Position = 0
				bot.SendMessage(message.Sender, "Hello, "+message.Sender.FirstName+"! I'm the daily Go bot ,\nChoose an option.", &telebot.SendOptions{
					ReplyMarkup: telebot.ReplyMarkup{
						Selective:       true,
						ForceReply:      true,
						CustomKeyboard:  keyboard,
						ResizeKeyboard:  true,
						OneTimeKeyboard: true,
					},
					ReplyTo: message,
				})
			} else if msg == "/cancel" {
				user.Position = 0
				bot.SendMessage(message.Sender, "Hello, "+message.Sender.FirstName+"! I'm the daily Go bot ,\nChoose an option.", &telebot.SendOptions{
					ReplyMarkup: telebot.ReplyMarkup{
						Selective:       true,
						ForceReply:      true,
						CustomKeyboard:  keyboard,
						ResizeKeyboard:  true,
						OneTimeKeyboard: true,
					},
					ReplyTo: message,
				})
			} else if msg == "New Link" {
				user.Position = 1
				bot.SendMessage(message.Sender, "Enter your link", nil)
			} else if msg == "New Post" {
				user.Position = 2
				bot.SendMessage(message.Sender, "Enter your post with markdown", nil)
			} else if msg == "Get Link Info" {
				user.Position = 3
				bot.SendMessage(message.Sender, "Enter the link ID", nil)
			} else if msg == "Post To Site" {
				user.Position = 4
				bot.SendMessage(message.Sender, "Enter the link ID", nil)
			} else {
				if user.Position == 1 {
					user.Position = 0
					id := newRandomID()
					for getlink(id).ID > 0 {
						id = newRandomID()
					}
					link := Link{}
					link.ID = id
					link.Link = strings.ToLower(msg)
					link.Save()

					bot.SendMessage(message.Sender, "Link is : \n\ngolng.ml/"+strconv.Itoa(id), &telebot.SendOptions{
						ReplyMarkup: telebot.ReplyMarkup{
							Selective:       true,
							ForceReply:      true,
							CustomKeyboard:  keyboard,
							ResizeKeyboard:  true,
							OneTimeKeyboard: true,
						},
						ReplyTo: message,
					})
				} else if user.Position == 2 {
					user.Position = 0
					bot.SendMessage(message.Sender, msg, &telebot.SendOptions{
						ReplyMarkup: telebot.ReplyMarkup{
							Selective:       true,
							ForceReply:      true,
							CustomKeyboard:  keyboard,
							ResizeKeyboard:  true,
							OneTimeKeyboard: true,
						},
						ReplyTo:   message,
						ParseMode: telebot.ModeMarkdown,
					})
				} else if user.Position == 3 {
					user.Position = 0
					id, _ := strconv.Atoi(msg)
					visits := getvisits(id)
					for _, visit := range visits {
						bot.SendMessage(message.Sender, visit.IP+" - "+strconv.FormatInt(visit.Date, 10), nil)
					}
					bot.SendMessage(message.Sender, "Total : "+strconv.Itoa(len(visits)), &telebot.SendOptions{
						ReplyMarkup: telebot.ReplyMarkup{
							Selective:       true,
							ForceReply:      true,
							CustomKeyboard:  keyboard,
							ResizeKeyboard:  true,
							OneTimeKeyboard: true,
						},
						ReplyTo:   message,
						ParseMode: telebot.ModeMarkdown,
					})
				} else if user.Position == 4 {
					id, _ := strconv.Atoi(msg)
					post := Post{}
					post.ID = id
					post.Save()
					user.Position = id
					bot.SendMessage(message.Sender, "Enter the message to post in site", &telebot.SendOptions{
						ReplyMarkup: telebot.ReplyMarkup{
							Selective:       true,
							ForceReply:      true,
							CustomKeyboard:  keyboard,
							ResizeKeyboard:  true,
							OneTimeKeyboard: true,
						},
						ReplyTo:   message,
						ParseMode: telebot.ModeMarkdown,
					})
				} else if user.Position > 0 { // new post to site , post id is in user.position
					post := getpost(user.Position)
					post.Content = msg
					post.Time = time.Now().Unix()
					post.Save()
					bot.SendMessage(message.Sender, "Message saved", &telebot.SendOptions{
						ReplyMarkup: telebot.ReplyMarkup{
							Selective:       true,
							ForceReply:      true,
							CustomKeyboard:  keyboard,
							ResizeKeyboard:  true,
							OneTimeKeyboard: true,
						},
						ReplyTo:   message,
						ParseMode: telebot.ModeMarkdown,
					})
					user.Position = 0
				}
			}
		} else {
			if msg == "/start" {
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
					err = bot.SendMessage(telebot.Chat{ID: ADMINID}, "Suggest : "+message.Text+" - @"+message.Sender.Username+" ["+message.Sender.FirstName+" "+message.Sender.LastName+"] ", nil)
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
					err = bot.SendMessage(telebot.Chat{ID: ADMINID}, "Contact : "+message.Text+" - @"+message.Sender.Username+" ["+message.Sender.FirstName+" "+message.Sender.LastName+"] ", nil)
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
		}

		user.Save()
	}

}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadFile("golng-static/index.html")
	w.Write(bytes)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<!DOCTYPE html>
		<html>
		<head>
		    <meta charset="utf-8">
		    <title>Golang</title>
			<style>
				body {
					background-color: #222;	
				}
				h1 {
					text-align:center;
					font-size:90px;
					color:#fff;	
				}
				@media only screen and (min-width:768px) {
					h1 {
						font-size:8vm;
					}
				}
			</style>
		</head>	
		<body>
			<h1>Not Found</h1>
		</body>
		</html>
	`))
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	errcode := 0
	posts := []Post{}

	time, err := strconv.Atoi(mux.Vars(r)["time"])
	if err != nil {
		errcode = 1
	} else {
		if time > 0 { // not first time
			db.Table("posts").Order("time desc").Where("time < ?", time).Limit(5).Scan(&posts)
		} else { // first time
			db.Table("posts").Order("time desc").Limit(5).Scan(&posts)
		}
	}
	bytes, _ := json.Marshal(&map[string]interface{}{
		"errorcode": errcode,
		"posts":     posts,
	})
	w.Write(bytes)
}

func ShortnerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.NotFound(w, r)
	}

	link := getlink(id)
	if link.ID != id {
		http.NotFound(w, r)
	}

	link.Count++
	link.Save()

	visit := Visit{
		Date: time.Now().Unix(),
		IP:   r.Header.Get("X-FORWARDED-FOR"),
		Link: id,
	}
	visit.Save()

	http.Redirect(w, r, link.Link, http.StatusFound)
}

func newRandomID() int {
	return 1000 + rand.Intn(9999-1000)
}
