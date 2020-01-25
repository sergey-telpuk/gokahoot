package models

type Player struct {
	ID     int    `gorm:"primary_key"`
	UUID   string `gorm:"unique_index"`
	Name   string `gorm:"not null"`
	Game   Game   `gorm:"foreignkey:GameID"`
	GameID int    `gorm:"type:integer REFERENCES games(id) ON DELETE CASCADE ON UPDATE CASCADE;not null"`
}
