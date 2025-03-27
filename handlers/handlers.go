package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"manga-catalog/database"
	"manga-catalog/models"
)

func GetMangaList(c *fiber.Ctx) error {
	fmt.Println("➡️ Вызван GetMangaList")

	var manga []models.Manga

	if database.DB == nil {
		fmt.Println("Ошибка: Подключение к базе отсутствует!")
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка соединения с базой данных"})
	}

	if err := database.DB.Find(&manga).Error; err != nil {
		fmt.Println("Ошибка запроса к базе:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка базы данных"})
	}

	fmt.Println("Найдено записей:", len(manga))
	return c.JSON(manga)
}

func CreateManga(c *fiber.Ctx) error {
	var manga models.Manga

	if err := c.BodyParser(&manga); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат JSON"})
	}

	if manga.Title == "" || manga.Description == "" || manga.Genre == "" || manga.Cover == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Все поля должны быть заполнены"})
	}

	if err := database.DB.Create(&manga).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка при добавлении манги"})
	}

	return c.Status(201).JSON(manga)
}
