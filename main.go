package main

import (
	"manga-catalog/database"
	"manga-catalog/handlers"
	"manga-catalog/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	user := r.Group("/api/user", middleware.AuthMiddleware())
	{
		user.GET("/profile", handlers.GetProfile)
		user.PUT("/profile", handlers.UpdateProfile)
	}

	admin := r.Group("/api/users", middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		admin.GET("/", handlers.GetAllUsers)
		admin.GET("/:id", handlers.GetUserByID)
		admin.POST("/", handlers.CreateUser)
		admin.PUT("/:id", handlers.UpdateUser)
		admin.DELETE("/:id", handlers.DeleteUser)
	}

	api := r.Group("/api", middleware.AuthMiddleware())
	{
		api.GET("/manga", handlers.GetMangaList)
		api.POST("/manga", handlers.CreateManga)
		api.GET("/manga/:id", handlers.GetMangaByID)
		api.PUT("/manga/:id", handlers.UpdateManga)
		api.DELETE("/manga/:id", handlers.DeleteManga)

		api.GET("/genres", handlers.GetAllGenres)
		api.GET("/genres/stats", handlers.GetGenresWithCount)

		api.POST("/manga/:id/comments", handlers.AddComment)
		api.GET("/manga/:id/comments", handlers.GetComments)
	}

	r.Run(":8080")
}
