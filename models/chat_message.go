package models

import "github.com/jinzhu/gorm"

type ChatMessage struct {
	gorm.Model
	UUID     string `gorm:"unique_index"`
	Message  string `gorm:"not null"`
	Game     Game
	GameID   int `gorm:"type:integer REFERENCES games(id) ON DELETE CASCADE ON UPDATE CASCADE;not null"`
	Player   Player
	PlayerID int `gorm:"type:integer REFERENCES players(id) ON DELETE CASCADE ON UPDATE CASCADE;not null"`
}
