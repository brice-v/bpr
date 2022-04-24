package main

import (
	"bpr/db"
	"bpr/handlers"
	"embed"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html"
)

//go:embed views/*
var viewsStaticDir embed.FS

//go:embed public/*
var publicStaticDir embed.FS

const DB_NAME = "bpr.db"

func main() {
	engine := html.NewFileSystem(http.FS(viewsStaticDir), ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	setupMiddleware(app)
	setupRoutes(app)

	app.Listen(":5961")
}

func setupMiddleware(app *fiber.App) {
	// use embedded public directory
	app.Use("/public", filesystem.New(filesystem.Config{
		Root:       http.FS(publicStaticDir),
		PathPrefix: "public",
	}))
	// use embedded favicon
	app.Use(favicon.New(favicon.Config{
		FileSystem: http.FS(publicStaticDir),
		File:       "public/favicon.ico",
	}))

	dbc, err := db.NewAndMigrate(DB_NAME)
	if err != nil {
		log.Fatal(err)
	}
	// store db connection in locals to be accessible by all reqs
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", dbc)
		return c.Next()
	})
	// 1 req/s
	app.Use(limiter.New(limiter.Config{
		Expiration: time.Second,
	}))
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
