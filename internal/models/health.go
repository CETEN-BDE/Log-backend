package models

import (
	"log-backend/autogen"

	"gorm.io/gorm"
)

type Health struct {
  gorm.Model  // Embeds ID, CreatedAt, UpdatedAt, DeletedAt
  autogen.Health
}