package database

import (
	"context"
	"fmt"
	"github.com/Niromash/niromash-api/api"
	"github.com/Niromash/niromash-api/utils/environment"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/nitishm/go-rejson/v4"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

var _ api.Database[api.RedisClient] = (*RedisDatabase)(nil)

type RedisDatabase struct {
	service api.MainService
	client  api.RedisClient
}

type redisClientImpl struct {
	client *redis.Client
	reJson *rejson.Handler
}

func (r *redisClientImpl) GetJsonString(key, path string) (string, error) {
	result, err := r.reJson.JSONGet(key, path)
	if err != nil {
		return "", err
	}
	str := string(result.([]byte))
	if str[0] != "\""[0] { // Check 2 first char (as rune) if it's not a quote
		return string(result.([]byte)), nil
	}
	return string(result.([]byte))[1 : len(str)-1], nil // [1:len(str)-1] - remove quotes
}

func (r *redisClientImpl) GetJsonInt(key, path string) (int, error) {
	result, err := r.reJson.JSONGet(key, path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(result.([]byte)))
}

func (r *redisClientImpl) GetJsonFloat(key, path string) (float64, error) {
	result, err := r.reJson.JSONGet(key, path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(result.([]byte)), 64)
}

func (r *redisClientImpl) GetJsonBool(key, path string) (bool, error) {
	result, err := r.reJson.JSONGet(key, path)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(string(result.([]byte)))
}

func (r *redisClientImpl) GetJsonObject(key, path string, object interface{}) error {
	if object == nil || (reflect.ValueOf(object).Kind() != reflect.Ptr) {
		return fmt.Errorf("object must be a pointer")
	}
	result, err := r.reJson.JSONGet(key, path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(result.([]byte), object); err != nil {
		return err
	}

	return nil
}

func (r *RedisDatabase) GetClient() api.RedisClient {
	return r.client
}

func (r *RedisDatabase) GetName() string {
	return "Redis"
}

func (r *RedisDatabase) Connect() (err error) {
	redisUrl, _ := url.Parse(environment.GetRedisUri())
	redisPassword, _ := redisUrl.User.Password()
	var redisClient redisClientImpl
	redisClient.client = redis.NewClient(&redis.Options{
		Addr:     redisUrl.Host,
		Password: redisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisClient.reJson = rejson.NewReJSONHandler()
	redisClient.reJson.SetGoRedisClient(redisClient.client)

	r.client = &redisClient

	return r.client.Base().Ping(ctx).Err()
}

func (r *RedisDatabase) Disconnect() error {
	return r.client.Base().Close()
}

func (r *RedisDatabase) Ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Base().Ping(ctx).Err() == nil
}

func (r *RedisDatabase) WaitForStart() {
	for {
		if r.client != nil && r.client.Base() != nil && r.Ping() {
			return
		}

		fmt.Printf("Waiting for %s to start...\n", r.GetName())
		time.Sleep(1 * time.Second)
	}
}

func (r *redisClientImpl) Base() *redis.Client {
	return r.client
}

func (r *redisClientImpl) ReJson() *rejson.Handler {
	return r.reJson
}
