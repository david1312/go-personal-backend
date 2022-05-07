package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//todo :
// 1. make migration command for database

func main(){
	fmt.Println("hai")
	// db, err := sql.Open("mysql", "sunmoris_dav:domainsemesta13@tcp(127.0.0.1)/sunmoris_customer?timeout=5s")
	db, err := sql.Open("mysql", "mysql_user:mysql_password@tcp(127.0.0.1:3306)/sunmoris_customer?timeout=30s")

    // if there is an error opening the connection, handle it
    if err != nil {
        fmt.Println(err.Error())
		return
    }
	// defer db.Close()
	err = db.Ping()
	if err != nil{
		fmt.Println("asd")
		fmt.Println(err.Error())
		return
	}

	// _, err = db.Query("INSERT INTO users VALUES ( 2, 'TEST' )")

    // // if there is an error inserting, handle it
    // if err != nil {
    //     panic(err.Error())
    // }
	// be careful deferring Queries if you are using transactions

	fmt.Println("sukses")
	//fmt.Println(db)
}