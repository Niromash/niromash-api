package api

import (
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

type RedisClient interface {
	Base() *redis.Client
	ReJson() *rejson.Handler
	GetJsonString(key, path string) (string, error)
	GetJsonInt(key, path string) (int, error)
	GetJsonFloat(key, path string) (float64, error)
	GetJsonBool(key, path string) (bool, error)
	GetJsonObject(key, path string, object interface{}) error
}
