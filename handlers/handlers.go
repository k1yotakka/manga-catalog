package handlers

import (
	"github.com/gin-gonic/gin"
	"manga-catalog/database"
	"manga-catalog/models"
	"net/http"
	"strconv"
)

func GetMangaList(c *gin.Context) {
	var manga []models.Manga

	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	genre := c.Query("genre")

	limit, err1 := strconv.Atoi(limitStr)
	page, err2 := strconv.Atoi(pageStr)
	if err1 != nil || err2 != nil || limit <= 0 || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные параметры пагинации"})
		return
	}
	offset := (page - 1) * limit

	query := database.DB.Limit(limit).Offset(offset)
	if genre != "" {
		query = query.Where("genre = ?", genre)
	}

	if err := query.Find(&manga).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  manga,
		"page":  page,
		"limit": limit,
	})
}

func CreateManga(c *gin.Context) {
	var manga models.Manga

	if err := c.ShouldBindJSON(&manga); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	if manga.Title == "" || manga.Description == "" || manga.Genre == "" || manga.Cover == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Все поля должны быть заполнены"})
		return
	}

	if err := database.DB.Create(&manga).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении манги"})
		return
	}

	c.JSON(http.StatusCreated, manga)
}

func GetMangaByID(c *gin.Context) {
	id := c.Param("id")
	var manga models.Manga

	if err := database.DB.First(&manga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Манга не найдена"})
		return
	}

	c.JSON(http.StatusOK, manga)
}

func UpdateManga(c *gin.Context) {
	id := c.Param("id")
	var manga models.Manga

	if err := database.DB.First(&manga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Манга не найдена"})
		return
	}

	var updatedData models.Manga
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	manga.Title = updatedData.Title
	manga.Description = updatedData.Description
	manga.Genre = updatedData.Genre
	manga.Cover = updatedData.Cover

	if err := database.DB.Save(&manga).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении манги"})
		return
	}

	c.JSON(http.StatusOK, manga)
}

func DeleteManga(c *gin.Context) {
	id := c.Param("id")
	var manga models.Manga

	if err := database.DB.First(&manga, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Манга не найдена"})
		return
	}

	if err := database.DB.Delete(&manga).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении манги"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Манга успешно удалена"})
}
