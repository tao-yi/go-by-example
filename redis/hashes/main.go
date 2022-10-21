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

type Model struct {
	Str1    string   `redis:"str1"`
	Str2    string   `redis:"str2"`
	Int     int      `redis:"int"`
	Bool    bool     `redis:"bool"`
	Ignored struct{} `redis:"-"`
}

func main() {
	rdb := client.Client()
	ctx := context.Background()
	key := "my-hash-model"
	_, err := rdb.HSet(ctx, key, map[string]interface{}{
		"str1": "hello",
		"str2": "world",
		"int":  123,
		"bool": 1,
	}).Result()
	exitOnError(err)

	var model1 Model
	if err = rdb.HGetAll(ctx, key).Scan(&model1); err != nil {
		panic(err)
	}
	fmt.Printf("%+v", model1)
}
