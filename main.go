package main

import (
	"log"
	"manga-catalog/database"
	"manga-catalog/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDB()

	app := fiber.New()

	app.Get("/manga", handlers.GetMangaList)
	app.Post("/manga", handlers.CreateManga)

	log.Fatal(app.Listen(":8080"))
}
