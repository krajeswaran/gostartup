package adapters

import (
	"errors"
	"github.com/go-redis/redis/v7"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/viper"
)

//CacheInterface interface for Cache adapter
type CacheInterface interface {
	DeepStatus() error
	GetApiStats() ([]string, error)
	UpdateApiStats(didApiFail bool) (int64, error)
	ResetApiStats() error
}

//CacheAdapter Struct to logically bind all the cache related functions
type CacheAdapter struct{
	cache redis.Cmdable
}

const (
	// StatCachePrefix Key prefix for stat related entries
	StatCachePrefix = "s:"
	//StatCacheHelloApiKey Key name for hello api stats
	StatCacheHelloApiKey = "hello_api"
	//StatCacheApiFailureCount Key name for counting API failures
	StatCacheApiFailureCount = "api_fail_count"
	//StatCacheApiTotalCount Key name for counting total API calls
	StatCacheApiTotalCount = "api_total_count"
)

//CacheInit initializes redis from config
func CacheInit() *CacheAdapter {
	c := &CacheAdapter{cache:nil}
	if viper.GetBool("USE_REDIS_CLUSTER") {
		cfg := redis.ClusterOptions{
			Addrs:        viper.GetStringSlice("cache_cluster_addresses"),
			Password:     viper.GetString("cache_password"),
			DialTimeout:  viper.GetDuration("cache_dial_timeout"),
			ReadTimeout:  viper.GetDuration("cache_read_timeout"),
			WriteTimeout: viper.GetDuration("cache_write_timeout"),
			PoolSize:     viper.GetInt("cache_pool_size"),
		}

		c.cache = redis.NewClusterClient(&cfg)
	} else {
		cfg := redis.Options{
			Addr:         viper.GetString("cache_address"),
			Password:     viper.GetString("cache_password"),
			DialTimeout:  viper.GetDuration("cache_dial_timeout"),
			ReadTimeout:  viper.GetDuration("cache_read_timeout"),
			WriteTimeout: viper.GetDuration("cache_write_timeout"),
			PoolSize:     viper.GetInt("cache_pool_size"),
		}

		c.cache = redis.NewClient(&cfg)
	}

	return c
}

//DeepStatus checks health of redis
func (c *CacheAdapter) DeepStatus() error {
	if err := c.cache.Ping().Err(); err != nil {
		return multierror.Append(err, errors.New("CACHE_REDIS_DOWN"))
	}
	return nil
}

//GetApiStats Retrieves API stats for hello service
func (c *CacheAdapter) GetApiStats() ([]string, error) {
	raw, err := c.cache.HMGet(StatCachePrefix+StatCacheHelloApiKey, StatCacheApiTotalCount, StatCacheApiFailureCount).Result()
	if err != nil {
		return nil, err
	}

	result := make([]string, len(raw))
	for i, v := range raw {
		result[i] = v.(string)
	}
	return result, nil
}

//UpdateApiStats Updates API stats in redis
func (c *CacheAdapter) UpdateApiStats(didApiFail bool) (int64, error) {
	count, err := c.cache.HIncrBy(StatCachePrefix+StatCacheHelloApiKey, StatCacheApiTotalCount, 1).Result()
	if err != nil {
		return 0, err
	}

	if didApiFail {
		_, err := c.cache.HIncrBy(StatCachePrefix+StatCacheHelloApiKey, StatCacheApiFailureCount, 1).Result()
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

//ResetApiStats Resets API stat counters
func (c *CacheAdapter) ResetApiStats() error {
	_, err := c.cache.Del(StatCachePrefix + StatCacheHelloApiKey).Result()
	if err != nil {
		return err
	}
	return nil
}
