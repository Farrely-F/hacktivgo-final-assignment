package utils

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var jwtKey []byte

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	jwtKey = []byte(os.Getenv("SECRET_JWT"))
}

func GenerateToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})

	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func ExtractUserIDFromToken(c *gin.Context) (uint, error) {
	// Get the JWT token from the request header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return 0, errors.New("missing token")
	}

	// Extract the token from the "Bearer <token>" format
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	} else {
		return 0, errors.New("invalid token format")
	}

	// Validate and parse the JWT token
	token, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	// Extract the user ID from the JWT claims
	claimsMap := token.Claims.(jwt.MapClaims)
	userID, ok := claimsMap["id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID in token")
	}

	return uint(userID), nil
}
