package bootstrap

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

func NewMysqlDB(dbUrl string) *sqlx.DB {
	db, err := sqlx.Connect("mysql", dbUrl)

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Fatalf("unable to connect to database : %v", err.Error())
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("unable to ping the database : %v", err.Error())
		return nil
	}

	db.SetConnMaxLifetime(time.Second * 30)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(100)
	return db
}
