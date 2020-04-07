package adapters

//DB Instance ref for DB adapter
var DB DBInterface

//Cache Instance ref for cache adapter
var Cache CacheInterface

// Init initializes "must-have" adapters - redis and db
func Init() {
	DB = DBInit()

	Cache = CacheInit()
}
