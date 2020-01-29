package models

type Answer struct {
	ID         int    `gorm:"primary_key"`
	Text       string `gorm:"not null"`
	Sequential int    `gorm:"not null"`
	ImgURL     *string
	Question   Question
	QuestionID int `gorm:"type:integer REFERENCES questions(id) ON DELETE CASCADE ON UPDATE CASCADE;not null"`
}
