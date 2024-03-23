package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username  string      		`json:"username" validate:"required,min=3,max=50" gorm:"uniqueIndex"`
    Email     string      		`json:"email" gorm:"uniqueIndex" validate:"required,email"`
    Password  string      		`json:"password" validate:"required,min=6"`
    Age       int         		`json:"age" validate:"required,gte=18"`
    CreatedAt time.Time			`json:"created_at"`
    UpdatedAt time.Time			`json:"updated_at"`
    Photos    []Photo			`json:"photos"`
    Comments  []Comment			`json:"comments"`
    SocialMedias []SocialMedia	`json:"social_medias"`
    ProfileImageURL string		`json:"profile_image_url" optional:"true"`
}

func (u *User) CustomValidationRules(v *validator.Validate) error {
    if u.Age < 18 {
        return fmt.Errorf("age must be greater than or equal to 18")
    }
    return nil
}