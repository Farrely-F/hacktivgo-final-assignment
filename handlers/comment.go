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


func GetComments(c *gin.Context) {
	var comments []models.Comment
	if err := db.DB.Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	for i, comment := range comments {
		var user models.User
		if err := db.DB.Where("id = ?", comment.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}
		comments[i].User = user
	}

	for i, comment := range comments {
		var photo models.Photo
		if err := db.DB.Where("id = ?", comment.PhotoID).First(&photo).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get photo"})
			return
		}
		comments[i].Photo = photo
	}

	commentResponses := make([]gin.H, len(comments))
	for i, comment := range comments {
		commentResponses[i] = gin.H{
			"id":        comment.ID,
			"comment":   comment.Message,
			"user_id":   comment.UserID,
			"created_at": comment.CreatedAt,
			"updated_at": comment.UpdatedAt,
			"user": gin.H{
				"id":     comment.UserID,
				"email":  comment.User.Email,
				"username": comment.User.Username,
			},
			"photo": gin.H{
				"id":        comment.PhotoID,
				"title":     comment.Photo.Title,
				"caption":   comment.Photo.Caption,
				"photo_url": comment.Photo.PhotoURL,
				"user_id":   comment.Photo.UserID,
			},
		}
	}

	c.JSON(http.StatusOK, commentResponses)
}

func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
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

	comment.UserID = userID

	if err := db.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	commentResponse := gin.H{
		"id": comment.ID,
		"message": comment.Message,
		"user_id": comment.UserID,
		"created_at": comment.CreatedAt,
	}

	c.JSON(http.StatusCreated, commentResponse)
}

func UpdateComment(c *gin.Context) {
	commentID := c.Param("commentId")
	var updatedComment models.Comment
	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := db.DB.Model(&models.Comment{}).Where("id = ?", commentID).Updates(updatedComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	var comment models.Comment
	if err := db.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated comment"})
		return
	}

	var photo models.Photo
	if err := db.DB.Where("id = ?", comment.PhotoID).First(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get photo"})
		return
	}

	var user models.User
	if err := db.DB.Where("id = ?", comment.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	updatedCommentResponse := gin.H{
		"id": commentID,
		"title": photo.Title,
		"caption": photo.Caption,
		"photo_url": photo.PhotoURL,
		"message": updatedComment.Message,
		"user_id": updatedComment.UserID,
		"updated_at": updatedComment.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedCommentResponse})
}

func DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")

	var comment models.Comment
	//Check and delete comment from database
	if err := db.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		}
		return
	}

	if err := db.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}