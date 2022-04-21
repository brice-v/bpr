package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"primaryKey"`
	Password string
}

func (u User) String() string {
	return fmt.Sprintf("User{Username: %s, Password: N/A}", u.Username)
}

type Post struct {
	gorm.Model
	Username  string
	Message   string
	Timestamp time.Time
}

func (p Post) String() string {
	return fmt.Sprintf("Post{Username: %s, Message: %s, Timestamp: %d}",
		p.Username, p.Message, p.Timestamp.Unix())
}

type Following struct {
	gorm.Model
	Username string
	Follower string
}

func (f Following) String() string {
	return fmt.Sprintf("Following{Username: %s, Follower: %s}", f.Username, f.Follower)
}

type Auth struct {
	gorm.Model
	Username string `gorm:"primaryKey"`
	AuthId   string
}

func (a Auth) String() string {
	return fmt.Sprintf("Auth{Username: %s, AuthId: %s}", a.Username, a.AuthId)
}
