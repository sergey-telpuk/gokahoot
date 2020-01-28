// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphql

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
