package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
)

type Db struct {
	con *gorm.DB
}

func (db *Db) GetConn() *gorm.DB {
	return db.con
}

func Init() *Db {
	con, err := gorm.Open("sqlite3", getFullPath()+"/gorm.db")

	if err != nil {
		log.Fatalf("Connection was failed %v", err)
	}

	return &Db{con: con}
}

func getFullPath() string {
	path, err := os.Getwd()

	if err != nil {
		log.Fatalf("%v", err)
	}

	return path
}
