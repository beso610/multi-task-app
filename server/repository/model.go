package repository

import (
	"time"

	"github.com/google/uuid"
)

// User ユーザー情報
type User struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	MailAddress string    `gorm:"type:varchar(256);not null"`
	Password    string    `gorm:"type:char(86);not null"`
	Salt        string    `gorm:"type:char(86);not null"`
	UserName    string    `gorm:"type:char(100);not null"`
	CreatedAt   time.Time `gorm:"precision:6"`
	UpdatedAt   time.Time `gorm:"precision:6"`
}

type Task struct {
	ID        uuid.UUID `gorm:"type:char(36);not null;primary_key"`
	UserID    uuid.UUID `gorm:"type:varchar(36);not null"`
	User      *User     `gorm:"foreignkey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ParentID  uuid.UUID `gorm:"type:varchar(36);not null"`
	Name      string    `gorm:"type:char(256);not null"`
	Finished  bool      `gorm:"type:bool;not null"`
	CreatedAt time.Time `gorm:"precision:6"`
	UpdatedAt time.Time `gorm:"precision:6"`
}
