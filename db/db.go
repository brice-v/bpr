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
	err = db.AutoMigrate(&models.Auth{})
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (d *DB) FindUser(username string) (*models.User, bool) {
	var user models.User
	txResult := d.db.First(&user, "username = ?", username)
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

func (d *DB) TotalUserCount() int {
	var count int64
	d.db.Table("users").Count(&count)
	return int(count)
}

func (d *DB) GetPosts(user *models.User) []models.Post {
	var posts []models.Post
	usersToGetFrom := []string{user.Username}
	var followings []models.Following
	result := d.db.Find(&followings, "follower = ?", user.Username)
	if result.RowsAffected != 0 {
		for _, v := range followings {
			usersToGetFrom = append(usersToGetFrom, v.Username)
		}
	}

	result = d.db.Find(&posts, "username in (?)", usersToGetFrom)
	if result.Error != nil {
		log.Printf("Do we need to handle this error? %s", result.Error.Error())
	}
	return posts
}

func (d *DB) GetSingleUsersPosts(user *models.User) []models.Post {
	var posts []models.Post
	result := d.db.Find(&posts, "username = ?", user.Username)
	if result.Error != nil {
		log.Printf("Handle this maybe? = %s", result.Error.Error())
	}
	return posts
}

func (d *DB) GetAllPosts() []models.Post {
	var posts []models.Post
	d.db.Find(&posts)
	return posts
}

func (d *DB) NewPost(p *models.Post) {
	tx := d.db.Create(p)
	if tx.Error != nil {
		log.Printf("NewPost failed: %s", tx.Error.Error())
	}
}

func (d *DB) GetAuthId(un string) (string, bool) {
	var a models.Auth
	result := d.db.First(&a, "username = ?", un)
	if result.Error != nil {
		return "", false
	}
	return a.AuthId, true
}

func (d *DB) SetAuthId(un, authId string) {
	_, ok := d.GetAuthId(un)
	if !ok {
		a := &models.Auth{
			Username: un,
			AuthId:   authId,
		}

		d.db.Create(a)
	} else {
		var a models.Auth
		d.db.First(&a, "username = ?", un)
		a.AuthId = authId
		d.db.Save(&a)
	}
}

func (d *DB) DeleteAuthId(un string) {
	var a models.Auth
	d.db.Where("username = ?", un).Delete(&a)
}

func (d *DB) IsUserFollowing(userToFollow, currentUser string) bool {
	if userToFollow == currentUser {
		return true
	}
	var f models.Following
	result := d.db.First(&f, "username = ? AND follower = ?", userToFollow, currentUser)
	if result.Error != nil {
		log.Printf("Should we handle this error? %s", result.Error.Error())
	}
	return result.RowsAffected == 1
}

func (d *DB) FollowUser(userToFollow, currentUser string) {
	if d.IsUserFollowing(userToFollow, currentUser) {
		return
	}
	f := models.Following{
		Username: userToFollow,
		Follower: currentUser,
	}
	result := d.db.Create(&f)
	if result.Error != nil {
		log.Printf("Should we handle this error? %s", result.Error.Error())
	}
}
