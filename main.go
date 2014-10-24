package main

import (
	"flag"
	"github.com/jinzhu/gorm"
	// _ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	crowdflow "github.com/workmachine/workmachine/crowdflow"
)

func init() {
	flag.StringVar(&crowdflow.Port, "port", "3000", "Port")
	flag.Parse()

	// Db, _ := gorm.Open("postgres", "user=gorm dbname=gorm sslmode=disable")
	Db, _ := gorm.Open("sqlite3", "gorm.db")
	crowdflow.Db = Db

	crowdflow.BuildDb()
}

func main() {
	crowdflow.StartServer()
}
