package main

import (
	"blog/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDatabase()
	database.CreateTables()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error { return c.JSON("On the air") })
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
