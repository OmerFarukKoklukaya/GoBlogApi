package main

import (
	"blog/database"
	"blog/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.ConnectDatabase()
	database.CreateTables()
	app := fiber.New()

	corsConfig := cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	})
	app.Use(corsConfig)

	app.Get("/", func(c *fiber.Ctx) error { return c.JSON("On the air") })
	router.Router(app)

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
