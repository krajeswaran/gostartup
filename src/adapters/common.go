package adapters

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Blank import for gorm-postgres
	"github.com/rs/zerolog/log"
	"gostartup/src/config"
	"time"
)

var db *gorm.DB
var cacheCluster *redis.ClusterClient
var cache *redis.Client
var DefaultCacheExpiry time.Duration

// Init for cache and db
func Init() {
	var err error

	// Initialize the redis server
	c := config.GetConfig()
	if c.GetBool("redis.cluster") {
		cacheCluster = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:       []string{c.GetString("redis.host_server")},
			Password:    "", // no password set
			ReadTimeout: time.Second,
		})

		if cacheCluster == nil {
			log.Error().Str("redis host", c.GetString("redis.host_server")).Msg("Can't connect to redis")
		}
	} else {
		cache = redis.NewClient(&redis.Options{
			Addr:        c.GetString("redis.host_server"),
			Password:    c.GetString("redis.password"),
			DB:          c.GetInt("redis.default_database"),
			ReadTimeout: time.Second,
		})
		if cache == nil {
			log.Error().Str("redis host", c.GetString("redis.host_server")).Msg("Can't connect to redis")
		}
	}
	DefaultCacheExpiry = time.Duration(config.GetConfig().GetInt("redis.default_key_expiration")) * time.Hour

	// create postgress connection
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s sslrootcert=%s connect_timeout=%s",
		c.GetString("postgres.db_host"), c.GetString("postgres.db_port"),
		c.GetString("postgres.db_user"), c.GetString("postgres.db_name"),
		c.GetString("postgres.db_ssl_mode"), c.GetString("postgres.db_password"),
		c.GetString("postgres.db_ssl_cert"), c.GetString("postgres.connect_timeout"))

	appName := c.GetString("postgres.application_name")
	if appName != "" {
		connectionString = connectionString + fmt.Sprintf(" application_name=%s", appName)
	}

	db, err = gorm.Open(c.GetString("postgres.db_type"), connectionString)
	if err != nil {
		panic("Can't connect to database, check config!" + err.Error())
	}

	connLifeTime := c.GetInt("postgres.conn_life_time")
	db.DB().SetConnMaxLifetime(time.Minute * time.Duration(connLifeTime))

	db.LogMode(c.GetBool("general.debug"))
}

// expose the db object
func Db() *gorm.DB {
	return db
}

// expose the cache object
func Cache() redis.Cmdable {
	c := config.GetConfig()
	if c.GetBool("redis.cluster") {
		return cacheCluster
	}
	return cache
}
