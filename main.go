package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer func(rdb *redis.Client) {
		err := rdb.Close()
		if err != nil {
			panic(err)
		}
	}(rdb)

	err := rdb.Set(ctx, "key", "hello, Redis!", 0).Err()
	if err != nil {
		log.Fatalf("Error setting key: %v", err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		log.Fatalf("Error getting key: %v", err)
	}
	fmt.Println("key:", val)

	val2, err := rdb.Get(ctx, "nonexistent").Result()
	if !errors.Is(err, redis.Nil) {
		if err != nil {
			log.Fatalf("Error getting nonexistent key: %v", err)
		} else {
			fmt.Println("nonexistent key:", val2)
		}
	} else {
		fmt.Println("key does not exist")
	}
}
