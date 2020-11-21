package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
)

func main() {
	app := fiber.New()
	app.Use(recover.New())

	app.Use(logger.New(logger.Config{
		Format:     "${time}|${pid}|${status}|${method}|${path}|${latency}\n",
		TimeFormat: "2006-01-02|15:04:05",
		TimeZone:   "America/Sao_Paulo",
	}))

	// Default middleware config
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	// Match any route
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("First handler")
		return c.Next()
	})

	// Match all routes starting with /api
	app.Use("/api", func(c *fiber.Ctx) error {
		log.Println("Second handler")
		panic("I'm an error")
		return c.Next()
	})

	// GET /api/register
	app.Get("/api/list", func(c *fiber.Ctx) error {
		log.Println("Last handler")
		return c.SendString("Hello, World 👋!")
	})

	// Starts the baile
	log.Fatal(app.Listen(":3000"))
}
