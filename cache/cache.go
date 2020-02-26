package cache

import (
	"github.com/go-redis/redis/v7"
)

const Prefix = "credo:metadata:v1:"
const KeyAuthWrite = Prefix + "auth:write"

var rdb *redis.Client

func connect() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

func init() {
	connect()
}

func CanWrite(token string) (bool, error) {
	return rdb.SIsMember(KeyAuthWrite, token).Result()
}
