package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

var ldb *LevelDB
var ldbFile = "acctTransDB"
var lonce sync.Once

type LevelDB struct {
	*leveldb.DB
}

func (l *LevelDB) KVDBName() string {
	return "LevelDB"
}

func GetLevelDBInstance() KeyValueDatabase {
	lonce.Do(func() {
		l, _ := leveldb.OpenFile(ldbFile, nil)
		ldb = &LevelDB{
			l,
		}
	})

	return ldb
}
