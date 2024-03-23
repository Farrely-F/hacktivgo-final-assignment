package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Photo struct {
    gorm.Model
    Title      string		`json:"title" validate:"required"`
    Caption    string		`json:"caption"`
    PhotoURL   string 		`json:"photo_url" validate:"required"`
    UserID     uint    		`json:"user_id"`
    User       User    		`json:"user" gorm:"foreignKey:UserID"` 
    CreatedAt  time.Time	`json:"created_at"`
    UpdatedAt  time.Time 	`json:"updated_at"`
    Comments   []Comment 	`json:"comments"` 
}

func (p *Photo) CustomValidationRules(v *validator.Validate) error {
    if len(p.Title) > 100 {
        return fmt.Errorf(p.Title, "title", "Title must be less than or equal to 100 characters")
    }
    return nil
}