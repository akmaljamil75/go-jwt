package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Username  string         `gorm:"not null;unique" json:"username"`
	Password  string         `gorm:"not null" json:"password"`
	RoleID    uint           `gorm:"not null" json:"role_id"`
	Version   int64          `gorm:"not null" json:"version"`
}
