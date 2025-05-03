package helper

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

func GetRedisCache(redisPool *redis.Pool, redisKey string) (string, error) {
	conn := redisPool.Get()
	defer conn.Close()

	cache, err := conn.Do("GET", redisKey)
	cacheValue, err := redis.String(cache, err)
	if err != nil {
		return "", err
	}
	return cacheValue, nil
}

func SetRedisCache(redisPool *redis.Pool, redisKey string, redisValue []byte, expireTime int32) error {
	conn := redisPool.Get()
	defer conn.Close()

	if expireTime < 1 {
		expireTime = 86400
	}

	_, errSet := conn.Do("SET", redisKey, string(redisValue)) // set the value to redis key
	if errSet != nil {
		log.Printf("error set redis cache with redisKey: %v, value: %v, errorMessage: %v", redisKey, string(redisValue), errSet)
		return errSet
	} else {
		_, errExpire := conn.Do("EXPIRE", redisKey, expireTime) // set cache available time to be 1 day long (86400 seconds), unless there is any updates to that user
		if errExpire != nil {
			log.Printf("error set redis cache expire time with redisKey: %v, value: %v, errorMessage: %v", redisKey, string(redisValue), errExpire)
			return errExpire
		}
	}
	return nil
}

func SaddRedisCache(redisPool *redis.Pool, redisSetKey string, value []byte) error {
	conn := redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", redisSetKey, value) // set the value to redis key
	if err != nil {
		log.Printf("error set redis cache with redisKey: %v, value: %v, errorMessage: %v", redisSetKey, value, err)
		return err
	}

	return nil
}

func SmembersRedisCache(redisPool *redis.Pool, redisSetKey string) ([]string, error) {
	conn := redisPool.Get()
	defer conn.Close()

	cache, err := conn.Do("SMEMBERS", redisSetKey)
	cacheValue, err := redis.Strings(cache, err)
	if err != nil {
		return nil, err
	}
	return cacheValue, nil
}

func DeleteRedisCache(redisPool *redis.Pool, redisKey string) error {
	conn := redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", redisKey) // delete the value of redis key
	if err != nil {
		log.Printf("error delete redis cache with redisKey: %v, errorMessage: %v", redisKey, err)
		return err
	}
	return nil
}

func GetIncrementValue(redisPool *redis.Pool, redisKey string, TTL ...int32) int32 {
	if redisKey == "" {
		return 0
	}

	conn := redisPool.Get()
	defer conn.Close()

	number, err := conn.Do("INCR", redisKey)

	// add expire time to the key
	if len(TTL) > 0 {
		_, errExpire := conn.Do("EXPIRE", redisKey, TTL[0])
		if errExpire != nil {
			log.Printf("error set redis cache expire time with redisKey: %v, errorMessage: %v", redisKey, errExpire)
		}
	}

	numberValue, err := redis.Int64(number, err)
	if err != nil {
		return 0
	}
	log.Print(numberValue)
	return int32(numberValue)
}

func GetOrSetCache[T any](redisPool *redis.Pool, key string, ttl time.Duration, enableCache bool, fetcher func() (T, error)) (T, error) {
	log.Printf("attempting to find cache with key: %v", key)

	var emptyResponse T

	if !enableCache {
		return fetcher()
	}

	conn := redisPool.Get()
	defer conn.Close()

	cacheStr, err := GetRedisCache(redisPool, key)
	if err == nil && cacheStr != "" {
		var response T
		if err := json.Unmarshal([]byte(cacheStr), &response); err == nil {
			log.Printf("Cache hit for key: %s", key)

			return response, nil
		}

		log.Printf("Failed to unmarshal cached value: %s", err)
	}

	// Cache is not available or error
	response, err := fetcher()
	if err != nil {
		return emptyResponse, err
	}

	cacheBytes, _ := json.Marshal(response)
	go func() {
		SetRedisCache(redisPool, key, cacheBytes, int32(ttl.Seconds()))
	}()

	return response, nil
}
