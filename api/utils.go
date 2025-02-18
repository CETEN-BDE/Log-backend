package api

import (
	"fmt"
	"log-backend/internal/models"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetAccountIDFromJWT(ctx echo.Context) (uint, error) {
	tokenRaw, err := ctx.Cookie("token")
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	token, err := jwt.Parse(tokenRaw.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		logrus.Error(err)
		return 0, err
	}
    claims := token.Claims.(jwt.MapClaims)
	accountID, ok := claims["AccountID"].(float64)
	if !ok {
		return 0, fmt.Errorf("AccountID not a uint in token claims")
	}
	return uint(accountID), nil
}

func IsMemberOfBDE(db *gorm.DB, accountID uint) (bool, error) {
	var account models.Account
	result := db.Where("id = ?", accountID).First(&account)
	if result.Error != nil {
		return false, result.Error
	}
	if account.ClubId == 1 {
		return true, nil
	}
	return false, nil
}