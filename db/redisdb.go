package db

import (
	_ "errors"
	"github.com/go-redis/redis"
	"sync"
)

var redb *RedisDB
var reAddr = "localhost:6379"
var redbPasswd = ""
var reonce sync.Once

type RedisDB struct {
	*redis.Client
}

func (r *RedisDB) KVDBName() string {
	return "RedisDB"
}

func GetRedisDBInstance() *RedisDB {
	reonce.Do(func() {
		r := redis.NewClient(&redis.Options{
			Addr:     reAddr,
			Password: redbPasswd,
			DB:       0,
		})

		redb = &RedisDB{
			r,
		}
	})

	return redb
}
