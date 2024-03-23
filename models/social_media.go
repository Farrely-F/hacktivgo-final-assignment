package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type SocialMedia struct {
    gorm.Model
    Name             string 	`json:"name" validate:"required,max=50"`
    SocialMediaURL   string 	`json:"social_media_url" validate:"required"`
    UserID           uint   	`json:"user_id"`
    User             User   	`json:"user" gorm:"foreignKey:UserID"`
    CreatedAt        time.Time	`json:"created_at"`
    UpdatedAt        time.Time	`json:"updated_at"`
}

func (s *SocialMedia) CustomValidationRules(v *validator.Validate) error {
    if len(s.Name) > 50 {
        return fmt.Errorf("name must be less than or equal to 50 characters")
    }
    return nil
}