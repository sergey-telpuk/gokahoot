package models

type Test struct {
	ID        int    `gorm:"primary_key"`
	UUID      string `gorm:"unique_index"`
	Name      string `gorm:"not null"`
	Questions []Question
}
