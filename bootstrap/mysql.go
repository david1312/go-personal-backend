package bootstrap

import (
	"database/sql"
	"log"
)


func NewMysqlDB(dbUrl string) *sql.DB {
	db, err := sql.Open("mysql",dbUrl)

    // if there is an error opening the connection, handle it
    if err != nil {
        log.Printf("unable to connect to database : %v", err.Error())
		return nil
    }

	err = db.Ping()
	if err != nil{
		log.Printf("unable to ping the database : %v", err.Error())
		return nil
	}
	return db
}