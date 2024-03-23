package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Comment struct {
    gorm.Model
    UserID     uint   		`json:"user_id" validate:"required"`
    User       User   		`json:"user" gorm:"foreignKey:UserID"`
    PhotoID    uint   		`json:"photo_id" validate:"required"`
    Photo      Photo  		`json:"photo" gorm:"foreignKey:PhotoID"`
    Message    string 		`json:"message" validate:"required"`
    CreatedAt  time.Time	`json:"created_at"`
    UpdatedAt  time.Time	`json:"updated_at"`
}

func (c *Comment) CustomValidationRules(v *validator.Validate) error {
    return nil
}