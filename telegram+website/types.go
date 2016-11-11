package main

type Link struct {
	ID    int `gorm:"primary_key AUTO_INCREMENT"`
	Link  string
	Count int
}

type Visit struct {
	ID   int `gorm:"primary_key AUTO_INCREMENT"`
	Link int
	IP   string
	Date int64
}

type User struct {
	UserId   int `gorm:"primary_key"`
	Position int
	LastUsed int64
}
