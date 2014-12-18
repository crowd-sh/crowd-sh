package main

import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	DatabaseUrlKey = "/workmachine/database_url"
)

var (
	Db gorm.DB
)

func dbConnect(databaseUrl string) {
	log.Println("Connecting to database:", databaseUrl)
	var err error
	Db, err = gorm.Open("postgres", databaseUrl)
	if err != nil {
		log.Println(err)
	}
	Db.LogMode(true)
}

func init() {
	etcdHosts := os.Getenv("ETCD_HOSTS")
	if etcdHosts == "" {
		etcdHosts = "http://127.0.0.1:4001"
	}

	etcdClient := etcd.NewClient([]string{etcdHosts})

	resp, err := etcdClient.Get(DatabaseUrlKey, false, false)
	if err != nil {
		panic(err)
	}

	databaseUrl := resp.Node.Value
	dbConnect(databaseUrl)
}

func main() {
	// go WorkExpirer()
	StartServer()
}
