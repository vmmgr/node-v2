package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"time"
)

const dbPath = "./node.db"

type VM struct {
	ID        int `gorm:primary_key`
	CreatedAt time.Time
	UpdatedAt time.Time
	GroupID   int
	Name      string
	CPU       int
	Mem       int
	Storage   string
	Net       string
	Status    int
	AutoStart bool
}

type Storage struct {
	ID        int `gorm:primary_key`
	CreatedAt time.Time
	UpdatedAt time.Time
	GroupID   int
	Name      string
	Driver    int //0:virtio
	Type      int //1:qcow2(default) 2:img
	Mode      int //0~9:AutoPath 10:ManualPath
	Path      string
	MaxSize   int
	Lock      int
}

type Net struct {
	ID         int `gorm:primary_key`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	GroupID    string
	Name       string
	Driver     int //0:virtio 1:e1000
	MacAddress string
	Vlan       int
	Status     int
}

type result struct {
	ID    int
	Error error
}

func InitCreateDB() {
	db := InitDB()
	defer db.Close()
}

func InitDB() *gorm.DB {
	return initSQLite3()
}

func initSQLite3() *gorm.DB {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("SQL open error")
	}
	//db.LogMode(true)
	db.SingularTable(true)

	db.AutoMigrate(&VM{})
	db.AutoMigrate(&Storage{})
	db.AutoMigrate(&Net{})

	return db
}
