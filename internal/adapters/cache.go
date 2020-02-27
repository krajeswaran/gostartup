package adapters

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/viper"
	"os"

	"github.com/go-redis/redis/v7"
)

//CacheAdapter Struct to logically bind all the cache related functions
type CacheAdapter struct{}

const (
	// StatCachePrefix Key prefix for stat related entries
	StatCachePrefix          = "s:"
	//StatCacheHelloApiKey Key name for hello api stats
	StatCacheHelloApiKey     = "hello_api"
	//StatCacheApiFailureCount Key name for counting API failures
	StatCacheApiFailureCount = "api_fail_count"
	//StatCacheApiTotalCount Key name for counting total API calls
	StatCacheApiTotalCount = "api_total_count"
)

//CacheInit initializes redis from env variables
func (c *CacheAdapter) CacheInit() (*redis.Cmdable, error) {
	var cache redis.Cmdable
	if viper.GetBool("USE_REDIS_CLUSTER") {
		var cfg redis.ClusterOptions
		if err:= viper.Unmarshal(&cfg); err != nil {
			return nil, multierror.Append(err, errors.New("unable to marshal redis config"))
		}

		cache = redis.NewClusterClient(&cfg)
	} else {
		var cfg redis.Options
		if err:= viper.Unmarshal(&cfg); err != nil {
			return nil, multierror.Append(err, errors.New("unable to marshal redis config"))
		}

		cache = redis.NewClient(&cfg)
	}

	if cache == nil {
		return nil, errors.New("redis host: " + os.Getenv("redis.host_server") + " Can't connect to redis")
	}

	return &cache, nil
}

//DeepStatus checks health of redis
func (c *CacheAdapter) DeepStatus() error {
	if err := cache.Ping().Err(); err != nil {
		return multierror.Append(err, errors.New("CACHE_REDIS_DOWN"))
	}
	return nil
}

//GetApiStats Retrieves API stats for hello service
func (c *CacheAdapter) GetApiStats() ([]string, error) {
	raw, err := cache.HMGet(StatCachePrefix+StatCacheHelloApiKey, StatCacheApiTotalCount, StatCacheApiFailureCount).Result()
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
	count, err := cache.HIncrBy(StatCachePrefix+StatCacheHelloApiKey, StatCacheApiTotalCount, 1).Result()
	if err != nil {
		return 0, err
	}

	if didApiFail {
		_, err := cache.HIncrBy(StatCachePrefix+StatCacheHelloApiKey, StatCacheApiFailureCount, 1).Result()
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

//ResetApiStats Resets API stat counters
func (c *CacheAdapter) ResetApiStats() error {
	_, err := cache.Del(StatCachePrefix+StatCacheHelloApiKey).Result()
	if err != nil {
		return err
	}
	return nil
}

