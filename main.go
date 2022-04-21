package main

import (
	"bpr/db"
	"bpr/handlers"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
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
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", dbc)
		return c.Next()
	})
	app.Use(limiter.New(limiter.Config{
		Expiration: time.Second,
	}))

	setupRoutes(app)

	app.Listen(":3000")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Index)
	app.Get("/signup", handlers.Signup)
	app.Get("/user/:username", handlers.User)
	app.Get("/monitor", monitor.New(monitor.Config{
		Next: monitorNextHelper,
	}))
	app.Get("/all", handlers.All)

	app.Post("/newUser", handlers.NewUser)
	app.Post("/login", handlers.Login)
	app.Post("/newPost", handlers.NewPost)
	app.Post("/logout", handlers.Logout)
	app.Post("/follow", handlers.Follow)
}

func monitorNextHelper(c *fiber.Ctx) bool {
	return !handlers.ValidateUser(c, "brice")
}
