package handlers

import (
	"bpr/models"
	"fmt"
	"log"
	"net/url"
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
			return c.Render("login", "")
		}
		if !validateUser(c, unc) {
			return c.Render("login", "")
		}
		// TODO: Render with all their messages
		posts := getDB(c).GetPosts(user)
		log.Printf("got to here, posts = %q", posts)
		return c.Render("home", fiber.Map{
			"Username": user.Username,
			"Posts":    posts,
		})
	}
	return c.Render("login", "")
}

// TODO: We can probably clean up the userAlreadyExists code
func Signup(c *fiber.Ctx) error {
	userAlreadyExistsQueryParam := c.Query("userAlreadyExists")
	if userAlreadyExistsQueryParam != "" {
		un, err := url.QueryUnescape(userAlreadyExistsQueryParam)
		if err != nil {
			return fmt.Errorf("failed to unescape '%s'. error: %s",
				userAlreadyExistsQueryParam, err.Error())
		}
		return c.Render("signup", fiber.Map{
			"UserAlreadyExists": un,
		})
	}
	return c.Render("signup", "")
}

func NewUser(c *fiber.Ctx) error {
	un := c.FormValue("username")
	pw := c.FormValue("password")
	log.Printf("un = %q, pw = %q", un, pw)
	dbc := getDB(c)
	if _, ok := dbc.FindUser(un); ok {
		redirect := fmt.Sprintf("/signup?userAlreadyExists=%s", url.QueryEscape(un))
		return c.Redirect(redirect)
	}
	dbc.CreateUser(un, pw)
	cache := getCache(c)
	authId := uuid.New().String()
	cache.Set(un, authId)
	setCookie(c, "username", un)
	setCookie(c, "authId", authId)
	return c.Redirect("/")
}

func Login(c *fiber.Ctx) error {
	un := c.FormValue("username")
	pw := c.FormValue("password")
	user, ok := getDB(c).FindUser(un)
	if !ok {
		return c.Render("login", fiber.Map{
			"Error": "User Not Found!",
		})
	}

	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pw)); err != nil {
		return c.Render("login", fiber.Map{
			"Error": "Invalid Credentials!",
		})
	}
	cache := getCache(c)
	authId := uuid.New().String()
	cache.Set(un, authId)
	setCookie(c, "username", un)
	setCookie(c, "authId", authId)
	return c.Redirect("/")
}

func User(c *fiber.Ctx) error {
	un := c.Params("username")
	// TODO: Just render the user and all their messages
	// If were not the same user then show a follow button
	return c.SendString("Hit User Endpoint with '" + un + "'")
}

func NewPost(c *fiber.Ctx) error {
	m := c.FormValue("message")
	un := c.FormValue("username")
	if !validateUser(c, un) {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}
	post := &models.Post{
		Message:        m,
		Timestamp:      time.Now(),
		AuthorUsername: un,
	}
	getDB(c).NewPost(post)
	return c.Redirect("/")
}
