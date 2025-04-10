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

func GetMangaByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var manga models.Manga

	if err := database.DB.First(&manga, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Манга не найдена"})
	}

	return c.JSON(manga)
}

func UpdateManga(c *fiber.Ctx) error {
	id := c.Params("id")
	var manga models.Manga

	if err := database.DB.First(&manga, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Манга не найдена"})
	}

	var updatedData models.Manga
	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Неверный формат JSON"})
	}

	manga.Title = updatedData.Title
	manga.Description = updatedData.Description
	manga.Genre = updatedData.Genre
	manga.Cover = updatedData.Cover

	if err := database.DB.Save(&manga).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка при обновлении манги"})
	}

	return c.JSON(manga)
}

func DeleteManga(c *fiber.Ctx) error {
	id := c.Params("id")
	var manga models.Manga

	if err := database.DB.First(&manga, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Манга не найдена"})
	}

	if err := database.DB.Delete(&manga).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Ошибка при удалении манги"})
	}

	return c.JSON(fiber.Map{"message": "Манга удалена успешно"})
}
