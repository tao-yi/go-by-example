package main

import (
	"context"
	"fmt"
	"go-by-example/redis/client"
	"time"
)

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	rdb := client.Client()
	ctx := context.Background()
	// 0 means no expiration time
	err := rdb.Set(ctx, "my-key", "my-value", 0).Err()
	exitOnError(err)

	exists, err := rdb.Exists(ctx, "my-key").Result()
	exitOnError(err)
	fmt.Println(exists)

	value := rdb.Get(ctx, "my-key").Val()
	fmt.Println(value)

	success, err := rdb.SetNX(ctx, "mykey", "lock", 1*time.Minute).Result()
	exitOnError(err)
	fmt.Println(success)

	// incr non-existing key will start from 0
	result, err := rdb.Incr(ctx, "not-existing-key").Result()
	exitOnError(err)
	fmt.Println(result)

	// MGET retrieves multiple string values in a single operation
	res, err := rdb.MSet(ctx, map[string]interface{}{
		"my-key1": "my-value1",
		"my-key2": "my-value2",
		"my-key3": "my-value3",
	}).Result()
	exitOnError(err)
	fmt.Println(res)

	values, err := rdb.MGet(ctx, "my-key1", "my-key2", "my-key3").Result()
	exitOnError(err)
	fmt.Println(values)
}
