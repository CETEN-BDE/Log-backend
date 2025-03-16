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
Password  string `gorm:"default:null"`
GoogleID  string `gorm:"unique;default:null"`
}
