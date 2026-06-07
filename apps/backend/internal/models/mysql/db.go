package mysql

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID           string `gorm:"primaryKey;size:64"`
	Email        string `gorm:"uniqueIndex;size:255;not null"`
	DisplayName  string `gorm:"size:128;not null"`
	PasswordHash string `gorm:"size:255"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func Open(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&GenerationProject{},
		&GenerationJob{},
		&GenerationChapter{},
		&GenerationScene{},
		&Prompt{},
	)
}
