package models

type Player struct {
	ID     int    `gorm:"primary_key"`
	UUID   string `gorm:"unique_index"`
	Name   string `gorm:"not null"`
	GameID int    `gorm:"type:integer REFERENCES games(id) ON DELETE CASCADE ON UPDATE CASCADE;not null"`
}
