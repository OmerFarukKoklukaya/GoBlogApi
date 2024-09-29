package main

import (
	"blog/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error { return c.JSON("On the air") })
	app.Listen(":3000")

	database.ConnectDatabase()
}
