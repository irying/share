package storage

import (
	"backend/conf"
	"backend/utils"
	"time"
	"github.com/go-redis/redis"
)

//RedisDriver redis client
type RedisDriver struct {
	Client *redis.Client
}

// Get get a key
func (r *RedisDriver) Get(key string) (interface{}, bool) {
	val, err := r.Client.Get(key).Result()
	if err == redis.Nil {
		return "", false
	} else if err != nil {
		panic(err)
	}
	return val, true
}

// Set set a key
func (r *RedisDriver) Set(key string, value interface{}, minutes time.Duration) {
	err := r.Client.Set(key, value, minutes).Err()
	if err != nil {
		panic(err)
	}
}

// Many mget many keys
func (r *RedisDriver) Many(keys []string) ([]interface{}) {
	val, err := r.Client.MGet(keys[:]...).Result()
	if err != nil {
		panic(err)
	}

	return val
}

// Put setnx a key
func (r *RedisDriver) Put(key string, value interface{}, minutes time.Duration) {
	r.Client.SetNX(key, value, minutes)
}

// PutMany put key during pipline
func (r *RedisDriver) PutMany(keys map[string]interface{}, minutes time.Duration) {
	pipe := r.Client.TxPipeline()
	for k, v := range keys {
		r.Put(k, v, minutes)
	}
	pipe.Exec()
}

// Increment Increment a key
func (r *RedisDriver) Increment(key string, value int) (int, bool) {

	val, err := r.Client.IncrBy(key, int64(value)).Result()
	if err != nil {
		panic(err)
	}
	return int(val), true
}

// Decrement Decrement a key
func (r *RedisDriver) Decrement(key string, value int) (int, bool) {
	panic("implement me")
}

// Forever a key forever exists
func (r *RedisDriver) Forever(key string, value interface{}) {
	panic("implement me")
}

// Forget remove a key
func (r *RedisDriver) Forget(key string) bool {
	panic("implement me")
}

// Flush remove all key
func (r *RedisDriver) Flush() bool {
	panic("implement me")
}

//NewRedis return redis client
func NewRedis() *RedisDriver {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Conf.Redis.Server,
		Password: conf.Conf.Redis.Pwd, // no password set
		DB:       0,                   // use default DB
		PoolSize: 10000,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		utils.ErrToPanic(err)
	}
	utils.LogPrint(pong)

	var redisDriver RedisDriver
	redisDriver.Client = client

	return &redisDriver
}
