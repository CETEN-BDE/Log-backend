package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Email     string `gorm:"unique"`
	FirstName string
	LastName  string
	Group     string
	ClubId    int 
	Password  string
}
