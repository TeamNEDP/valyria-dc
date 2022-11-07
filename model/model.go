package model

import "gorm.io/gorm"

type User struct {
	ID     string `gorm:"primaryKey"`
	Name   string `gorm:"uniqueIndex"`
	Email  string `gorm:"uniqueIndex"`
	Avatar []byte
	Rating int `gorm:"default:0"`
}

type UserScript struct {
	gorm.Model
	UserID string
	User   User
	Name   string
	Code   string
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&UserScript{},
	)
}
