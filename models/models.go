package models

import (
	"fmt"
	"time"
)

type User struct {
	Username string `gorm:"primaryKey"`
	Password string
}

func (u User) String() string {
	return fmt.Sprintf("User{Username: %s, Password: N/A}", u.Username)
}

type Post struct {
	Message        string
	Timestamp      time.Time
	AuthorUsername string `gorm:"primaryKey"`
}

func (p Post) String() string {
	return fmt.Sprintf("Post{AuthorUsername: %s, Message: %s, Timestamp: %d}",
		p.AuthorUsername, p.Message, p.Timestamp.Unix())
}

type Following struct {
	Username string `gorm:"primaryKey"`
	Follower string
}

func (f Following) String() string {
	return fmt.Sprintf("Following{Username: %s, Follower: %s}", f.Username, f.Follower)
}
