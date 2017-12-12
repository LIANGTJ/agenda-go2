package model

import (
	"github.com/jinzhu/gorm"

	_ "github.com/mattn/go-sqlite3"

	"util"
)

var agendaDB *gorm.DB

func init() {
	// db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")  // from gorm's doc
	db, err := gorm.Open("sqlite3", "/tmp/gorm.db")
	util.PanicIf(err)
	agendaDB = db
}
