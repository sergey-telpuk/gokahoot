package models

type Test struct {
	ID        int        `gorm:"primary_key"`
	UUID      string     `gorm:"unique_index"`
	Code      string     `gorm:"not null"`
	Name      string     `gorm:"not null"`
	Questions []Question `gorm:"ForeignKey:TestID"`
}
