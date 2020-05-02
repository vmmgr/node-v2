package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const dbName = "./node.db"

type VM struct {
	ID          int
	Name        string
	CPU         int
	Mem         int
	Storage     string
	StoragePath string
	Net         string
	Vnc         int
	Socket      string
	Status      int
	AutoStart   bool
}

func connectdb() *sql.DB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Println(err)
		fmt.Println("SQL open error")
	}

	return db
}

func Createdb() bool {
	db := *connectdb()
	defer db.Close()

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "vm" ("id" INTEGER PRIMARY KEY, "name" VARCHAR(255), "cpu" INT,"memory" INT, "storage" VARCHAR(500),"storagepath" VARCHAR(500),"net" VARCHAR(255),"vnc" INT, "socket" VARCHAR(255),"status" INT,"autostart" boolean)`)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error: SQL open Error")
		return false
	}
	return true
}
