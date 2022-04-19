package db

import (
	"bpr/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

func NewAndMigrate(dbName string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(dbName))
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.Post{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.Following{})
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (d *DB) FindUser(username string) (*models.User, bool) {
	var user models.User
	txResult := d.db.First(&user, "username = ?", username)
	log.Printf("FindUser txResult = %d", txResult.RowsAffected)
	// TODO: Probably want to return the user if we find them
	if txResult.RowsAffected == 0 {
		return nil, false
	}
	return &user, true
}

func (d *DB) CreateUser(username, password string) {
	// Generate "hash" to store from user password
	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
	}
	d.db.Create(&models.User{
		Username: username,
		Password: string(pwHash),
	})
}

func (d *DB) GetPosts(user *models.User) []models.Post {
	var posts []models.Post
	usersToGetFrom := []string{user.Username}
	var followings []models.Following
	result := d.db.Find(&followings, user.Username)
	if result.RowsAffected != 0 {
		for _, v := range followings {
			usersToGetFrom = append(usersToGetFrom, v.Follower)
		}
	}

	result = d.db.Find(&posts, usersToGetFrom)
	if result.Error != nil {
		log.Printf("Do we need to handle this error? %s", result.Error.Error())
	}
	return posts
}
