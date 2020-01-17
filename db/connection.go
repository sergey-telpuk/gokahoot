package db

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
	"time"
)

const ContainerName = "db"

type Db struct {
	con *gorm.DB
}

func (db *Db) GetConn() *gorm.DB {
	return db.con
}

func Init() (*Db, error) {
	con, err := gorm.Open("sqlite3", getFullPath()+"/gorm.db")

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	con.DB().SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	con.DB().SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	con.DB().SetConnMaxLifetime(time.Hour)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Connection was failed %v", err))
	}

	return &Db{con: con}, nil
}

func getFullPath() string {
	path, err := os.Getwd()

	if err != nil {
		log.Fatalf("%v", err)
	}

	return path
}
