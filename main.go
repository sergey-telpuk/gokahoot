package main

import (
	"encoding/json"
	guuid "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/sergey-telpuk/gokahoot/bot"
	"github.com/sergey-telpuk/gokahoot/db"
	"github.com/sergey-telpuk/gokahoot/models"
	"github.com/sergey-telpuk/gokahoot/server"
	"github.com/sergey-telpuk/gokahoot/services"
	"gopkg.in/gormigrate.v1"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	_ = godotenv.Load()
	s, err := services.New()

	if err != nil {
		log.Fatal(err)
	}

	go bot.Init(s).Run()
	migrate(s)
	createTest(s)

	if err := server.Run(s); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)

	defer s.Clean()
}

func migrate(di *services.DI) {
	sDB := di.Container.Get(db.ContainerName).(*db.Db)

	sDB.GetConn().AutoMigrate(
		&models.Test{},
		&models.Question{},
		&models.Answer{},
		&models.Answer{},
		&models.Game{},
		&models.Player{},
		&models.PlayerAnswer{},
		&models.ChatMessage{},
	)
	log.Printf("Migration did run successfully")
}

type verb struct {
	word         string
	descriptions []string
}

func createTest(di *services.DI) {
	sDB := di.Container.Get(db.ContainerName).(*db.Db)

	verbs := map[string]*verb{}

	var jsonData map[string]map[string]interface{}

	dir, _ := os.Getwd()

	file, err := ioutil.ReadFile(dir + "/fixtures/phrasal.verbs.json")

	if err != nil {
		log.Fatalf("Reading file was failed: %v", err)
	}

	if err := json.Unmarshal([]byte(file), &jsonData); err != nil {
		log.Fatalf("Parsing json was failed: %v", err)
	}
	it := 0
	for word, v := range jsonData {
		it++
		var descriptions []string
		s, ok := v["descriptions"].([]interface{})
		if !ok {
			continue
		}
		for _, d := range s {
			descriptions = append(descriptions, d.(string))
		}

		verbs[word] = &verb{
			word:         word,
			descriptions: descriptions,
		}

		if it > 35 {
			break
		}
	}

	m := gormigrate.New(sDB.GetConn(), gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: guuid.New().String(),
			Migrate: func(tx *gorm.DB) error {
				test := &models.Test{
					UUID: guuid.New().String(),
					Name: "Phrasal Verbs",
				}
				if err := tx.Create(test).Error; err != nil {
					return err
				}

				for _, v := range verbs {
					rightNumber := random(0)
					question := &models.Question{
						UUID:        guuid.New().String(),
						Text:        v.word,
						TestID:      test.ID,
						RightAnswer: rightNumber,
					}

					if err := tx.Create(question).Error; err != nil {
						return err
					}
					m := &models.Answer{
						Text:       strings.Join(v.descriptions, ","),
						Sequential: rightNumber,
						QuestionID: question.ID,
					}
					if err := tx.Create(m).Error; err != nil {
						return err
					}
					for i := 0; i < 3; i++ {
						_verb := randomVerb(v.word, verbs)
						m := &models.Answer{
							Text:       strings.Join(_verb.descriptions, ","),
							Sequential: random(rightNumber),
							QuestionID: question.ID,
						}
						if err := tx.Create(m).Error; err != nil {
							return err
						}
					}
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration did run successfully")
}

func randomVerb(exclude string, verbs map[string]*verb) *verb {
	verbsArray := []string{}
	for word, _ := range verbs {
		verbsArray = append(verbsArray, word)
	}
	rand.Seed(time.Now().Unix())
	min := 0
	max := len(verbsArray)
	n := rand.Intn(max-min) + min

	if verbsArray[n] == exclude {
		return randomVerb(exclude, verbs)
	}
	return verbs[verbsArray[n]]
}

//TODO
func random(exclude int) int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 1000
	n := rand.Intn(max-min) + min
	if n == exclude {
		return random(exclude)
	}
	return n
}
