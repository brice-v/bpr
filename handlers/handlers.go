package handlers

import (
	"bpr/models"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Index(c *fiber.Ctx) error {
	unc := c.Cookies("username")
	if unc != "" {
		user, userFound := getDB(c).FindUser(unc)
		if !userFound {
			return c.Render("views/login", "")
		}
		if !ValidateUser(c, unc) {
			return c.Render("views/login", "")
		}
		posts := getDB(c).GetPosts(user)
		return c.Render("views/home", fiber.Map{
			"Username": user.Username,
			"Posts":    posts,
		})
	}
	return c.Render("views/login", "")
}

func Signup(c *fiber.Ctx) error {
	return c.Render("views/signup", "")
}

func NewUser(c *fiber.Ctx) error {
	un := c.FormValue("username")
	pw := c.FormValue("password")
	dbc := getDB(c)
	if _, ok := dbc.FindUser(un); ok {
		errorString := fmt.Sprintf("Username '%s' already exists!", un)
		return c.Render("views/signup", fiber.Map{
			"Error": errorString,
		})
	}
	if dbc.TotalUserCount() > 100 {
		return c.SendString("Max users of 100 reached")
	}
	dbc.CreateUser(un, pw)
	authId := uuid.New().String()
	getDB(c).SetAuthId(un, authId)
	setCookie(c, "username", un)
	setCookie(c, "authId", authId)
	return c.Redirect("/")
}

func Login(c *fiber.Ctx) error {
	un := c.FormValue("username")
	pw := c.FormValue("password")
	user, ok := getDB(c).FindUser(un)
	if !ok {
		return c.Render("views/login", fiber.Map{
			"Error": "User Not Found!",
		})
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pw)); err != nil {
		return c.Render("views/login", fiber.Map{
			"Error": "Invalid Credentials!",
		})
	}
	authId := uuid.New().String()
	getDB(c).SetAuthId(un, authId)
	setCookie(c, "username", un)
	setCookie(c, "authId", authId)
	return c.Redirect("/")
}

func User(c *fiber.Ctx) error {
	// username were getting to
	un := c.Params("username")
	currentUn := c.Cookies("username")
	if !ValidateUser(c, currentUn) {
		return c.Redirect("/")
	}
	dbc := getDB(c)
	user, ok := dbc.FindUser(un)
	if !ok {
		return c.SendString("User '" + un + "' not found!")
	}
	posts := dbc.GetSingleUsersPosts(user)
	return c.Render("views/user", fiber.Map{
		"Username":           un,
		"IsUsernameLoggedIn": un == currentUn,
		"FollowerUsername":   currentUn,
		"Posts":              posts,
		"IsNotFollowing":     !dbc.IsUserFollowing(un, currentUn),
	})
}

func NewPost(c *fiber.Ctx) error {
	m := c.FormValue("message")
	un := c.FormValue("username")
	if !ValidateUser(c, un) {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}
	post := &models.Post{
		Message:   m,
		Timestamp: time.Now(),
		Username:  un,
	}
	getDB(c).NewPost(post)
	return c.Redirect("/")
}

func Logout(c *fiber.Ctx) error {
	un := c.FormValue("username")
	c.ClearCookie("username", "authId")
	getDB(c).DeleteAuthId(un)
	return c.Redirect("/")
}

func Follow(c *fiber.Ctx) error {
	un := c.FormValue("username")
	f := c.FormValue("follower")
	getDB(c).FollowUser(un, f)
	return c.Redirect("/user/" + un)
}

func All(c *fiber.Ctx) error {
	posts := getDB(c).GetAllPosts()
	return c.Render("views/all", fiber.Map{
		"Posts": posts,
	})
}
