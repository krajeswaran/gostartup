package adapters

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-redis/redis/v7"
)

var db *pg.DB
var cache redis.Cmdable

//DBInstance Instance ref for DB adapter
var DBInstance = DBAdapter{}

//CacheInstance Instance ref for cache adapter
var CacheInstance = CacheAdapter{}

// Init initializes "must-have" adapters - in our case redis and db
func Init() {
	db = DBInstance.DBInit()

	cachePtr := CacheInstance.CacheInit()
	cache = *cachePtr
}
