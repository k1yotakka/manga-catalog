package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"manga-catalog/database"
	"manga-catalog/handlers"
	"manga-catalog/middleware"
	"manga-catalog/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func init() {
	os.Setenv("DB_URL", "postgres://postgres:2705@localhost:5432/manga_test?sslmode=disable")
	database.ConnectDB()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	return r
}

func TestRegisterSuccess(t *testing.T) {
	router := setupRouter()

	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	body := map[string]string{"username": username, "password": "testpass"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestRegisterInvalidFormat(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestLoginInvalidFormat(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestLoginWrongCredentials(t *testing.T) {
	router := setupRouter()

	body := map[string]string{"username": "wronguser", "password": "wrongpass"}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestLoginSuccess(t *testing.T) {
	router := setupRouter()

	body := map[string]string{"username": "loginuser", "password": "test123"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Now test login
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetNonExistentManga(t *testing.T) {
	router := gin.Default()
	router.GET("/manga/:id", handlers.GetMangaByID)

	req, _ := http.NewRequest("GET", "/manga/999999", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestCreateMangaInvalid(t *testing.T) {
	router := gin.Default()
	router.POST("/manga", handlers.CreateManga)

	body := `{"title": "", "description": "", "genre": "", "cover": ""}`
	req, _ := http.NewRequest("POST", "/manga", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestCreateMangaValid(t *testing.T) {
	router := gin.Default()
	router.POST("/manga", handlers.CreateManga)

	body := `{"title": "Bleach", "description": "Soul Reapers", "genre": "Action", "cover": "cover.jpg"}`
	req, _ := http.NewRequest("POST", "/manga", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestGetMangaByIDNotFound(t *testing.T) {
	router := gin.Default()
	router.GET("/manga/:id", handlers.GetMangaByID)

	req, _ := http.NewRequest("GET", "/manga/9999", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestUpdateMangaNotFound(t *testing.T) {
	router := gin.Default()
	router.PUT("/manga/:id", handlers.UpdateManga)

	body := `{"title": "Updated", "description": "Desc", "genre": "Drama", "cover": "cover.jpg"}`
	req, _ := http.NewRequest("PUT", "/manga/9999", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func getTokenAndMangaID(t *testing.T) (string, uint) {
	router := setupRouter()

	username := fmt.Sprintf("favuser_%d", time.Now().UnixNano())
	password := "pass123"
	creds := map[string]string{"username": username, "password": password}
	body, _ := json.Marshal(creds)

	// Регистрация
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Логин
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var loginResp map[string]string
	json.Unmarshal(resp.Body.Bytes(), &loginResp)
	token := loginResp["token"]

	// Добавляем мангу в БД напрямую
	manga := models.Manga{
		Title:       "Test Favorite Manga",
		Description: "Fav desc",
		Genre:       "FavGenre",
		Cover:       "cover.jpg",
	}
	database.DB.Create(&manga)

	return token, manga.ID
}

func setupFavoriteRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/manga/:id/favorite", handlers.AddToFavorites)
		api.GET("/favorites", handlers.GetFavorites)
	}
	return r
}

func TestAddGetFavorite(t *testing.T) {
	token, mangaID := getTokenAndMangaID(t)
	router := setupFavoriteRouter()

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/manga/%d/favorite", mangaID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	req, _ = http.NewRequest("GET", "/api/favorites", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string][]models.Manga
	_ = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.GreaterOrEqual(t, len(result["favorites"]), 1)

}
