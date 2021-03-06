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

type BroadcasChatGame struct {
	UUID    string           `json:"UUID"`
	Message string           `json:"message"`
	Player  *BroadcastPlayer `json:"player"`
	Time    string           `json:"time"`
}

type BroadcastAnswerForChartGame struct {
	AnswerID int                             `json:"answerID"`
	Players  []*BroadcastPlayersForChartGame `json:"players"`
}

type BroadcastPlayer struct {
	UUID     string `json:"UUID"`
	GameCode string `json:"gameCode"`
	Name     string `json:"name"`
}

type BroadcastPlayersForChartGame struct {
	Player   *BroadcastPlayer `json:"player"`
	WasRight bool             `json:"wasRight"`
}

type BroadcastPlayingGame struct {
	CurrentTimeSec      int                            `json:"currentTimeSec"`
	GameCode            string                         `json:"gameCode"`
	StartTimeSec        int                            `json:"startTimeSec"`
	CurrentQuestionUUID string                         `json:"currentQuestionUUID"`
	GameStatusEnum      GameStatus                     `json:"gameStatusEnum"`
	Answers             []*BroadcastAnswerForChartGame `json:"answers"`
}

type ChatMessage struct {
	UUID    string  `json:"UUID"`
	Message string  `json:"message"`
	Player  *Player `json:"player"`
	Game    *Game   `json:"game"`
	Time    string  `json:"time"`
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

type ReportAnswer struct {
	Answer *Answer `json:"answer"`
	Right  bool    `json:"right"`
}

type ReportGame struct {
	Code    string          `json:"code"`
	Players []*ReportPlayer `json:"players"`
}

type ReportPlayer struct {
	Player  *Player         `json:"player"`
	Answers []*ReportAnswer `json:"answers"`
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

type ChatTimeOrder string

const (
	ChatTimeOrderDesc ChatTimeOrder = "DESC"
	ChatTimeOrderAsc  ChatTimeOrder = "ASC"
)

var AllChatTimeOrder = []ChatTimeOrder{
	ChatTimeOrderDesc,
	ChatTimeOrderAsc,
}

func (e ChatTimeOrder) IsValid() bool {
	switch e {
	case ChatTimeOrderDesc, ChatTimeOrderAsc:
		return true
	}
	return false
}

func (e ChatTimeOrder) String() string {
	return string(e)
}

func (e *ChatTimeOrder) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ChatTimeOrder(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ChatTimeOrder", str)
	}
	return nil
}

func (e ChatTimeOrder) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
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
