package models

type Game struct {
	ID     int    `gorm:"primary_key"`
	Code   string `gorm:"unique_index"`
	TestID int    `gorm:"type:integer REFERENCES tests(id) ON DELETE CASCADE ON UPDATE CASCADE;not null"`
}
