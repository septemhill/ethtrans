package db

import (
	"sync"

	"github.com/go-pg/pg"
)

var ronce sync.Once
var rdb *pg.DB

const username = "septemhill"
const password = "gintamaed3op2"
const dbname = "ether"

func GetRDBInstance() *pg.DB {
	ronce.Do(func() {
		rdb = pg.Connect(&pg.Options{
			User:     username,
			Password: password,
			Database: dbname,
		})
	})

	return rdb
}
