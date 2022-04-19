package models

import "time"

type User struct {
	Username string `gorm:"primaryKey"`
	Password string
}

type Post struct {
	Message        string
	Timestamp      time.Time
	AuthorUsername string
}

type Following struct {
	Username string `gorm:"primaryKey"`
	Follower string
}
