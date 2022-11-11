package model

import (
	"gorm.io/gorm"
	"time"
	"valyria-dc/game"
)

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

type Game struct {
	ID        string `gorm:"uniqueIndex"`
	Finished  bool   `gorm:"default:false"`
	RScriptID uint
	RScript   UserScript `gorm:"foreignKey:RScriptID"`
	BScriptID uint
	BScript   UserScript `gorm:"foreignKey:BScriptID"`
	Setting   game.GameSetting
	Ticks     game.GameTicks
	Result    game.GameResult
	CreatedAt time.Time
}

type UserCompetition struct {
	UserID       uint `gorm:"uniqueIndex:idx_compete_user_script"`
	User         User
	UserScriptID uint `gorm:"uniqueIndex:idx_compete_user_script"`
	UserScript   UserScript
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&UserScript{},
		&UserSession{},
		&Game{},
		&UserCompetition{},
	)
}
