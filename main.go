package main

import (
	"hacktiv-go/final-project-test/db"
	"hacktiv-go/final-project-test/handlers"
	"hacktiv-go/final-project-test/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()

	r := gin.Default()

	r.POST("/users/register", handlers.RegisterUser)
	r.POST("/users/login", handlers.LoginUser)
	r.PUT("/users/:userId", middleware.AuthMiddleware(), handlers.UpdateUser)
	r.DELETE("/users", middleware.AuthMiddleware(), handlers.DeleteUser)

	r.POST("/photos", middleware.AuthMiddleware(), handlers.CreatePhoto)
	r.GET("/photos", middleware.AuthMiddleware(), handlers.GetPhotos)
	r.PUT("/photos/:photoId", middleware.AuthMiddleware(), handlers.UpdatePhoto)
	r.DELETE("/photos/:photoId", middleware.AuthMiddleware(), handlers.DeletePhoto)

	r.POST("/comments", middleware.AuthMiddleware(), handlers.CreateComment)
	r.GET("/comments", middleware.AuthMiddleware(), handlers.GetComments)
	r.PUT("/comments/:commentId", middleware.AuthMiddleware(), handlers.UpdateComment)
	r.DELETE("/comments/:commentId", middleware.AuthMiddleware(), handlers.DeleteComment)

	r.POST("/socialmedias", middleware.AuthMiddleware(), handlers.CreateSocialMedia)
	r.GET("/socialmedias", middleware.AuthMiddleware(), handlers.GetSocialMedias)
	r.PUT("/socialmedias/:socialMediaId", middleware.AuthMiddleware(), handlers.UpdateSocialMedia)
	r.DELETE("/socialmedias/:socialMediaId", middleware.AuthMiddleware(), handlers.DeleteSocialMedia)

	r.Run(":8080") // Sesuaikan dengan port yang Anda inginkan
}
