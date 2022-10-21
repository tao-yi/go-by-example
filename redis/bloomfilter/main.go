package main

import (
	"context"
	"fmt"
	"go-by-example/redis/client"
)

func failOnError(msg string, err error) {
	if err != nil {
		panic(fmt.Errorf("%s: %v", msg, err.Error()))
	}
}

func main() {
	cli := client.Client()
	ctx := context.Background()
	key := "bf_key"
	val := "value1"
	inserted, err := cli.Do(ctx, "BF.ADD", key, val).Bool()
	failOnError("add key failed", err)
	if inserted {
		fmt.Println("value1 inserted")
	}

	exists, err := cli.Do(ctx, "BF.EXISTS", key, val).Bool()
	failOnError("BF.EXISTS error", err)
	if exists {
		fmt.Printf("%s %s exists\n", key, val)
	}
}
