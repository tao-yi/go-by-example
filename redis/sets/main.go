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
	key := "user:123"
	rdb.Del(ctx, key)

	// LPUSH add to head
	_, err := rdb.SAdd(ctx, key, 101, 101, 359, 404).Result()
	exitOnError(err)

	// RPush add to tail
	_, err = rdb.SRem(ctx, key, 101).Result()
	exitOnError(err)

	// test for set membership
	exist, err := rdb.SIsMember(ctx, key, 101).Result()
	exitOnError(err)
	fmt.Println(exist)
}
