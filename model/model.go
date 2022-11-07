package model

import "gorm.io/gorm"

type User struct {
	ID       string `gorm:"primaryKey"`
	Name     string `gorm:"uniqueIndex"`
	Email    string `gorm:"uniqueIndex"`
	Avatar   []byte
	Rating   int    `gorm:"default:0"`
	Password []byte `json:"password"`
}

type UserScript struct {
	gorm.Model
	UserID string `gorm:"uniqueIndex:idx_user_script_name"`
	User   User
	Name   string `gorm:"uniqueIndex:idx_user_script_name"`
	Code   string
}

type UserSession struct {
	gorm.Model
	UserID    string
	User      User
	UserAgent string
	SessionID string `gorm:"uniqueIndex"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&UserScript{},
		&UserSession{},
	)
}
