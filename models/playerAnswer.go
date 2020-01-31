package models

type PlayerAnswer struct {
	ID         int `gorm:"primary_key"`
	Player     Player
	PlayerID   int `gorm:"type:integer REFERENCES players(id) ON DELETE CASCADE ON UPDATE CASCADE;not null"`
	Game       Game
	GameID     int `gorm:"type:integer REFERENCES games(id) ON DELETE CASCADE ON UPDATE CASCADE;not null"`
	Question   Question
	QuestionID int `gorm:"type:integer REFERENCES questions(id) ON DELETE CASCADE ON UPDATE CASCADE;not null"`
	Answer     Answer
	AnswerID   int `gorm:"type:integer REFERENCES answers(id) ON DELETE CASCADE ON UPDATE CASCADE;not null"`
	WasRight   bool
}
