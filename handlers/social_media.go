package handlers

import (
	"errors"
	"fmt"
	"hacktiv-go/final-project-test/db"
	"hacktiv-go/final-project-test/models"
	"hacktiv-go/final-project-test/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func GetSocialMedias(c *gin.Context) {
	var socialMedias []models.SocialMedia
	if err := db.DB.Find(&socialMedias).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get social medias"})
		return
	}

	for i, socialMedia := range socialMedias {
		var user models.User
		if err := db.DB.Where("id = ?", socialMedia.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}
		socialMedias[i].User = user
	}

	socmedResponse := make([]gin.H, len(socialMedias))
	for i, socialMedia := range socialMedias {
		socmedResponse[i] = gin.H{
			"id": socialMedia.ID,
			"name": socialMedia.Name,
			"social_media_url": socialMedia.SocialMediaURL,
			"created_at": socialMedia.CreatedAt,
			"user": gin.H{
				"id": socialMedia.UserID,
				"profile_image_url": socialMedia.User.ProfileImageURL,
				"username": socialMedia.User.Username,
			},
		}
	}

	c.JSON(http.StatusOK, gin.H{"social_medias": socmedResponse})
}

func CreateSocialMedia(c *gin.Context) {
	var socialMedia models.SocialMedia
	if err := c.ShouldBindJSON(&socialMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID, err := utils.ExtractUserIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	socialMedia.UserID = userID
	if err := db.DB.Create(&socialMedia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	socialMediaResponse := gin.H{
		"id": socialMedia.ID,
		"name": socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user_id": socialMedia.UserID,
		"created_at": socialMedia.CreatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{"message": socialMediaResponse})
}


func UpdateSocialMedia(c *gin.Context) {
	socialMediaID := c.Param("socialMediaId")
	var updatedSocialMedia models.SocialMedia
	if err := c.ShouldBindJSON(&updatedSocialMedia); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	socialMedia := models.SocialMedia{}
	if err := db.DB.Where("id = ?", socialMediaID).First(&socialMedia).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get social media"})
		return
	}

	if err := db.DB.Model(&socialMedia).Updates(updatedSocialMedia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update social media: %s", err)})
		return
	}

	socialMediaResponse := gin.H{
		"id": socialMedia.ID,
		"name": socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user_id": socialMedia.UserID,
		"updated_at": socialMedia.UpdatedAt,
	}

	c.JSON(http.StatusOK, socialMediaResponse)
}

func DeleteSocialMedia(c *gin.Context) {
	socialMediaID := c.Param("socialMediaId")

	socialMedia := models.SocialMedia{}
	if err := db.DB.Where("id = ?", socialMediaID).First(&socialMedia).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get social media"})
		return
	}

	if err := db.DB.Delete(&socialMedia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete social media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Your social media has been succesfully deleted"})
}