package main

import (
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/jwt/v2"
	"log"
	"time"
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

	// Static files
	app.Static("/", "./public")

	//// Favicon
	//app.Use(favicon.New(favicon.Config{
	//	File: "./public/favicon.ico",
	//}))

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

	// Home is always acessible
	app.Get("/", home)

	// Login route
	app.Post("/login", login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	// Restricted Routes
	app.Get("/restricted", restricted)

	// Starts the baile
	log.Fatal(app.Listen(":3000"))
}

// Para servir o html, renomeie o arquivo /public/index.html1 para .html
func home(c *fiber.Ctx) error {
	log.Printf("Carregando home handler...")
	return c.SendString("Home, sweet home!")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"]
	email := claims["email"]
	return c.SendString(fmt.Sprintf("Welcome, %s (%s)!", name, email))
}

// curl --data "user=diego&pass=123" http://localhost:3000/login
// curl localhost:3000/restricted -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZW1haWwiOiJ0ZXN0QHRlc3QuY29tIiwiZXhwIjoxNjA2MDM5MTAxLCJuYW1lIjoiRGllZ28gRmluY2F0dG8ifQ.16KRqos7CBjZCAqL0ERJuI5NzN0_sRzHPjQbWSO5cgY"
func login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Throws Unauthorized error
	if user != "diego" || pass != "123" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Diego Fincatto"
	claims["email"] = "test@test.com"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 6).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}
