package api

import (
	"log-backend/internal/models"

	"gorm.io/gorm"
)

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
