package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client          // declare redis client variable/pointer to redis connection
	Ctx = context.Background() //cretes background connections for redis to support timeouts/cancellation
)

// func connect to redis server
func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{ //create new redis client and assign to it to redis
		Addr:     "localhost:6379", // default redis port
		Password: "",               //no password set
		DB:       0,                //default database
	})

	//ping the redis server to check connection
	status, err := RDB.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to redis: %v", err))
	}
	fmt.Printf("Connected to redis! status: %s\n", status)
}
