package bootstrap

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
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
	return db
}
