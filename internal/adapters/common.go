package adapters

import (
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var cache redis.Cmdable

//DBInstance Instance ref for DB adapter
var DBInstance = DBAdapter{}

//CacheInstance Instance ref for cache adapter
var CacheInstance = CacheAdapter{}

// Init initializes "must-have" adapters - in our case redis and db
func Init() {
	var err error

	db, err = DBInstance.DBInit()
	if err != nil {
		panic("Can't connect to database, check config: " + err.Error())
	}

	var cachePtr *redis.Cmdable
	cachePtr, err = CacheInstance.CacheInit()
	if err != nil {
		panic("Can't connect to cache, check config: " + err.Error())
	}
	cache = *cachePtr
}
