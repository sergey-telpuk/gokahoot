package models

const (
	GameInWaitingPlayers = iota //default
	GameInPlaying
	GameInFinished
)

type Game struct {
	ID     int    `gorm:"primary_key:true"`
	Code   string `gorm:"unique_index"`
	Test   *Test  `gorm:"foreignkey:TestID"`
	TestID int    `gorm:"type:integer REFERENCES tests(id) ON DELETE CASCADE ON UPDATE CASCADE;not null"`
	Status int    `gorm:"type:integer;default:0;not null"`
}
