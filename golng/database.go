package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var err error
var db *gorm.DB

func open_db() {
	db, err = gorm.Open("sqlite3", "database.sqlite3")
	if err != nil {
		panic(err)
	}

	if !db.HasTable(&Link{}) {
		db.CreateTable(&Link{})
	}

	if !db.HasTable(&Visit{}) {
		db.CreateTable(&Visit{})
	}

	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
	}

	if !db.HasTable(&Post{}) {
		db.CreateTable(&Post{})
	}
}

func close_db() {
	err = db.Close()
	if err != nil {
		panic(err)
	}
}

func getlink(a int) (link Link) {
	db.Table("links").Find(&link, "id = ?", a)
	return
}

func getpost(a int) (post Post) {
	db.Table("posts").Find(&post, "id = ?", a)
	return
}

func (link Link) Save() {
	db.Table("links").Save(&link)
}

func (visit Visit) Save() {
	db.Table("visits").Save(&visit)
}

func (user User) Save() {
	db.Table("users").Save(&user)
}

func (post Post) Save() {
	db.Table("posts").Save(&post)
}

func getvisits(a int) (visits []Visit) {
	db.Table("visits").Order("date", true).Where("link = ?", a).Scan(&visits)
	return
}
