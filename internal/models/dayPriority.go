package models

import (
	"gorm.io/gorm"
)

type DayPriority struct {
	gorm.Model
	AccountID uint
	Day       int
	Priority  int
}