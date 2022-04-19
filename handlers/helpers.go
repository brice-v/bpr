package handlers

import (
	"bpr/db"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func getDB(c *fiber.Ctx) *db.DB {
	dbc, ok := c.Locals("db").(*db.DB)
	if !ok {
		log.Fatal("Database Connection not found in Locals")
	}
	return dbc
}

func getCache(c *fiber.Ctx) *db.Cache {
	cache, ok := c.Locals("cache").(*db.Cache)
	if !ok {
		log.Fatal("Cache not found in Locals")
	}
	return cache
}

func setCookie(c *fiber.Ctx, key, value string) {
	cook := new(fiber.Cookie)
	cook.Name = key
	cook.Value = value
	cook.Expires = time.Now().Add(time.Hour)
	c.Cookie(cook)
}
