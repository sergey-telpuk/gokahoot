package db

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"time"
)

const (
	ContainerName   = "db"
	DefaultDriverDb = "sqlite3"
)

type Db struct {
	con *gorm.DB
}

func (db *Db) GetConn() *gorm.DB {
	return db.con
}

func Init() (*Db, error) {
	var con *gorm.DB
	var err error
	driverDb := os.Getenv("DRIVER_DB")
	if driverDb == "" {
		driverDb = DefaultDriverDb
	}

	switch DefaultDriverDb {
	case "sqlite3":
		con, err = sqlite3()
	default:
		con, err = gorm.Open("postgres", "host=localhost port=5433 user=gokahoot dbname=gokahoot password=gokahoot sslmode=disable")
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	con.DB().SetMaxIdleConns(100)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	con.DB().SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	con.DB().SetConnMaxLifetime(time.Hour)

	con.LogMode(true)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Connection was failed %v", err))
	}

	return &Db{con: con}, nil
}

func sqlite3() (*gorm.DB, error) {
	return gorm.Open("sqlite3", getFullPath()+"/gorm.db")
}

func getFullPath() string {
	path, err := os.Getwd()

	if err != nil {
		log.Fatalf("%v", err)
	}

	return path
}
