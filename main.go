package main

import (
	"bpr/db"
	"bpr/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

const DB_NAME = "test.db"

func main() {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	app.Static("/", "./public")

	dbc, err := db.NewAndMigrate(DB_NAME)
	if err != nil {
		log.Fatal(err)
	}
	cache := db.NewCache()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", dbc)
		c.Locals("cache", cache)
		return c.Next()
	})

	app.Get("/", handlers.Index)
	app.Get("/signup", handlers.Signup)
	app.Get("/user/:username", handlers.User)
	app.Post("/newUser", handlers.NewUser)
	app.Post("/login", handlers.Login)
	app.Post("/newPost", handlers.NewPost)

	app.Listen(":3000")
}
