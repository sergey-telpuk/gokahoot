package models

type Question struct {
	ID      int      `gorm:"primary_key"`
	UUID    string   `gorm:"unique_index"`
	TestID  int      `gorm:"type:integer REFERENCES tests(id) ON DELETE CASCADE ON UPDATE CASCADE;not null"`
	Answers []Answer `gorm:"ForeignKey:QuestionID"`

	Text        string `gorm:"not null"`
	ImgURL      *string
	RightAnswer int `gorm:"not null"`
}
