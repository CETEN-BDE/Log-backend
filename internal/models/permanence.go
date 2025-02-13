package models

import (
	"log-backend/autogen"
	"time"
	"gorm.io/gorm"
)

type Permanence struct {
	gorm.Model
	Date time.Time
	Status autogen.PermanenceStatus
	AccountID uint
	Account Account
}