package main

import (
	"github.com/jinzhu/gorm"
)

var (
	Port string
	Db   gorm.DB
)

func BuildDb() {
	Db.AutoMigrate(&User{})
	Db.AutoMigrate(&Workflow{})
}
