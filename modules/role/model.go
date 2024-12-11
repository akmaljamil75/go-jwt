package role

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name      string         `gorm:"not null;unique" json:"name"`
	Version   int64          `gorm:"not null" json:"version"`
}
