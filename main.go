package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: zerolog.NewConsoleWriter().Out})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer func(rdb *redis.Client) {
		err := rdb.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Error closing Redis connection")
		}
	}(rdb)

	err := rdb.Set(ctx, "key", "hello, Redis!", 0).Err()
	if err != nil {
		log.Fatal().Err(err).Msg("Error setting key in Redis")
	} else {
		log.Info().Msg("Successfully set key in Redis")
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting key from Redis")
	} else {
		log.Info().Str("key", "key").Str("value", val).
			Msg("Successfully retrieved key from Redis")
	}
	fmt.Println("key:", val)

	val2, err := rdb.Get(ctx, "nonexistent").Result()
	if !errors.Is(err, redis.Nil) {
		if err != nil {
			log.Fatal().Err(err).Msg("Error getting nonexistent key from Redis")
		} else {
			log.Info().Str("key", "nonexistent").Str("value", val2).
				Msg("Nonexistent key retrieved")
		}
	} else {
		log.Warn().Msg("Key does not exist in Redis")
		fmt.Println("key does not exist")
	}
}
