package handlers

import (
	"errors"
	"hacktiv-go/final-project-test/db"
	"hacktiv-go/final-project-test/models"
	"hacktiv-go/final-project-test/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPhotos(c *gin.Context) {
	var photos models.Photo
	if err := db.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, username, email")
    }).Find(&photos).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get photos"})
        return
    }

	photoResponse := gin.H{
		"id":         	photos.ID,
		"title":      	photos.Title,
		"caption":    	photos.Caption,
		"photo_url":  	photos.PhotoURL,
		"user_id":    	photos.UserID,
		"created_at": 	photos.CreatedAt,
		"updated_at": 	photos.UpdatedAt,
		"user": 		gin.H{
			"email": 	photos.User.Email,
			"username": photos.User.Username,
		},
	}

	c.JSON(http.StatusOK, photoResponse)
}

func CreatePhoto(c *gin.Context) {
	var photo models.Photo
	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID, err := utils.ExtractUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	photo.UserID = userID

	if err := db.DB.Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create photo"})
		return
	}

	if err := db.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, email")
	}).Find(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get photo"})
		return
	}

	createdPhotoResponse := gin.H{
		"id": photo.ID,
		"title": photo.Title,
		"caption": photo.Caption,
		"photo_url": photo.PhotoURL,
		"user_id": photo.UserID,
		"created_at": photo.CreatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{"data": createdPhotoResponse})
}




func UpdatePhoto(c *gin.Context) {
	photoID := c.Param("photoId")
	var updatedPhoto models.Photo

	if err := db.DB.Where("id = ?", photoID).First(&updatedPhoto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update photo"})
		}
		return
	}

	if err := c.ShouldBindJSON(&updatedPhoto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := db.DB.Model(&models.Photo{}).Where("id = ?", photoID).Updates(updatedPhoto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update photo"})
		return
	}

	updatedPhotoResponse := gin.H{
		"id": updatedPhoto.ID,
		"title": updatedPhoto.Title,
		"caption": updatedPhoto.Caption,
		"photo_url": updatedPhoto.PhotoURL,
		"user_id": updatedPhoto.UserID,
		"updated_at": updatedPhoto.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedPhotoResponse})
}

func DeletePhoto(c *gin.Context) {
	photoID := c.Param("photoId")

	var photo models.Photo
	if err := db.DB.Where("id = ?", photoID).First(&photo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
		}
		return
	}

	if err := db.DB.Delete(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your photo has been deleted succesfully"})
}
