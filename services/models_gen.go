// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package services

import (
	"fmt"
	"io"
	"strconv"
)

type Answer struct {
	ID         int     `json:"ID"`
	Text       string  `json:"text"`
	Sequential int     `json:"sequential"`
	ImgURL     *string `json:"imgURL"`
}

type BroadcastPlayer struct {
	UUID     string `json:"UUID"`
	GameCode string `json:"gameCode"`
	Name     string `json:"name"`
}

type BroadcastPlayingGame struct {
	Timer               int        `json:"timer"`
	GameCode            string     `json:"gameCode"`
	CurrentQuestionUUID string     `json:"currentQuestionUUID"`
	GameStatusEnum      GameStatus `json:"gameStatusEnum"`
}

type Game struct {
	Test    *Test     `json:"test"`
	Code    string    `json:"CODE"`
	Status  int       `json:"Status"`
	Players []*Player `json:"players"`
}

type InputAnswer struct {
	Sequential int     `json:"sequential"`
	Text       string  `json:"text"`
	ImgURL     *string `json:"imgURL"`
}

type InputJoinPlayer struct {
	GameCode string `json:"gameCode"`
	Name     string `json:"name"`
}

type NewQuestion struct {
	TestUUID    string         `json:"testUUID"`
	Text        string         `json:"text"`
	ImgURL      *string        `json:"imgURL"`
	RightAnswer int            `json:"rightAnswer"`
	Answers     []*InputAnswer `json:"answers"`
}

type NewTest struct {
	Name string `json:"name"`
}

type Player struct {
	UUID string `json:"UUID"`
	Game *Game  `json:"game"`
	Name string `json:"name"`
}

type Question struct {
	ID          int       `json:"ID"`
	UUID        string    `json:"UUID"`
	TestID      int       `json:"testID"`
	Text        string    `json:"text"`
	ImgURL      *string   `json:"imgURL"`
	RightAnswer int       `json:"rightAnswer"`
	Answers     []*Answer `json:"answers"`
}

type StartGame struct {
	GameCode string `json:"gameCode"`
}

type Status struct {
	Success bool `json:"success"`
}

type Test struct {
	ID        int         `json:"ID"`
	UUID      string      `json:"UUID"`
	Name      string      `json:"name"`
	Questions []*Question `json:"questions"`
}

type UpdateAnswer struct {
	ID         int     `json:"ID"`
	Sequential *int    `json:"sequential"`
	Text       *string `json:"text"`
	ImgURL     *string `json:"imgURL"`
}

type UpdateQuestion struct {
	UUID        string          `json:"UUID"`
	Text        *string         `json:"text"`
	ImgURL      *string         `json:"imgURL"`
	RightAnswer *int            `json:"rightAnswer"`
	Answers     []*UpdateAnswer `json:"answers"`
}

type UpdateTest struct {
	UUID      string            `json:"UUID"`
	Name      string            `json:"name"`
	Questions []*UpdateQuestion `json:"questions"`
}

type GameStatus string

const (
	GameStatusWaitForPlayers GameStatus = "WAIT_FOR_PLAYERS"
	GameStatusPlaying        GameStatus = "PLAYING"
	GameStatusFinished       GameStatus = "FINISHED"
)

var AllGameStatus = []GameStatus{
	GameStatusWaitForPlayers,
	GameStatusPlaying,
	GameStatusFinished,
}

func (e GameStatus) IsValid() bool {
	switch e {
	case GameStatusWaitForPlayers, GameStatusPlaying, GameStatusFinished:
		return true
	}
	return false
}

func (e GameStatus) String() string {
	return string(e)
}

func (e *GameStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GameStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GameStatus", str)
	}
	return nil
}

func (e GameStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
