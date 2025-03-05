package infrastucture

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "othello_user:pass@tcp(localhost:3306)/othello?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	return db
}
