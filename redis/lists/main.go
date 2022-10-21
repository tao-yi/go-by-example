package main

import (
	"context"
	"fmt"
	"go-by-example/redis/client"
)

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}

/**
Redis lists are linked lists of string values
*/

func main() {
	rdb := client.Client()
	ctx := context.Background()
	key := "work:queue:ids"
	rdb.Del(ctx, key)

	// LPUSH add to head
	_, err := rdb.LPush(ctx, key, 101, 237, 359, 404).Result()
	exitOnError(err)

	// RPush add to tail
	_, err = rdb.RPush(ctx, key, 1, 2, 3, 4, 5).Result()
	exitOnError(err)

	// LLen returns the length of a list
	count := rdb.LLen(ctx, key).Val()
	fmt.Println(count)

	val, err := rdb.LPop(ctx, key).Int()
	exitOnError(err)
	fmt.Println(val)

	val, err = rdb.LPop(ctx, key).Int()
	exitOnError(err)
	fmt.Println(val)
}
